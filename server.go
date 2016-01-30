package yar

import (
	"bytes"
	"net"
	"reflect"
	"yar/packager"
	"yar/transports"
	"errors"
	"fmt"
)

type Server struct {
	handler_list map[string]Handler

	opt          map[string]interface{}

	transport    transports.Transport
}

var CONNECTION_TOTAL = 0

func NewServer(host string, port int) *Server {

	s := new(Server)

	s.handler_list = make(map[string]Handler)

	tran, err := transports.NewTcp(host, port)

	if err != nil {
		panic(err)
	}

	s.transport = tran
	return s

}

func (self *Server) RegisterHandler(name string, handler Handler) bool {

	self.handler_list[name] = handler

	return true
}

func (self *Server) RemoveHandler(name string) bool {

	delete(self.handler_list, name)

	return true
}

func (self *Server) ReadProtocol(conn net.Conn) (*Protocol, error) {

	protocol := NewProtocol()
	protocol_len := PROTOCOL_LENGTH
	protocol_buffer := make([]byte, protocol_len)
	read_total := 0

	for {

		size, err := conn.Read(protocol_buffer[read_total:])

		if err != nil {
			return nil, err
		}

		read_total += size

		if read_total >= int(protocol_len) {
			break
		}

	}

	buffer := bytes.NewBuffer(protocol_buffer)
	protocol.Init(buffer)
	protocol.BodyLength -= 8

	return protocol, nil
}

func (self *Server) getProtocol(conn net.Conn) (protocol *Protocol, err error) {

	protocol_buffer := make([]byte, PROTOCOL_LENGTH)
	receive_len := 0

	for {

		_len, err := conn.Read(protocol_buffer[receive_len:])

		if err != nil {
			return nil, err
		}

		receive_len += _len

		if receive_len >= PROTOCOL_LENGTH {
			break
		}
	}

	protocol = NewProtocolWithBytes(bytes.NewBuffer(protocol_buffer))
	//todo check data
	return protocol, nil
}

func (self *Server) getRequest(conn net.Conn,protocol *Protocol) (request *Request, err error) {

	if protocol == nil {

		return nil,errors.New("get or parse protocol failed")

	}

	//min body length is 8 . packager is 8 bit
	if protocol.BodyLength < 8 {

		return nil,errors.New("body data errror")

	}

	real_body_length := protocol.BodyLength - 8

	body_buffer := make([]byte,real_body_length)
	receive_len := 0

	for {
		_len,err := conn.Read(body_buffer[receive_len:])
		if err != nil {
			return nil,err
		}
		receive_len += _len
		if receive_len >=int(real_body_length) {
			break
		}
	}

	request = new(Request)
	request.Protocol = protocol
	err = packager.Unpack(protocol.Packager[:],body_buffer,&request)
	return request,err
}

func (self *Server) sendResponse(conn net.Conn,response *Response) {

	protocol := response.Protocol

	if protocol == nil {
		conn.Close()
		return
	}

	ret, _ := packager.Pack(protocol.Packager[0:], response)
	protocol.BodyLength = uint32(len(ret) + 8)
	send_protocol := protocol.Bytes()
	conn.Write(send_protocol.Bytes())
	conn.Write(ret)
	conn.Close()
}

func (self *Server) handler(conn net.Conn) {

	var handler Handler = nil
	var fv reflect.Value
	var call_params []interface{}
	var real_params []reflect.Value
	var rs []reflect.Value

	var protocol *Protocol
	var request *Request
	response := NewResponse()

	protocol,err := self.getProtocol(conn)

	response.Protocol = protocol

	if err != nil {
		conn.Close()
		return
	}

	defer func() {

		if err := recover(); err != nil {
			response.Status = ERR_EMPTY_RESPONSE
			response.Error = fmt.Sprintf("%s:%s", "server has panic!",err)
			self.sendResponse(conn,response)
		}

	}()

	request,err = self.getRequest(conn,protocol)

	if err != nil {

		response.Status = ERR_REQUEST
		response.Error = "server read request error"
		self.sendResponse(conn,response)
		return
	}

	request.Protocol = protocol
	request.Id = request.Protocol.Id
	handler = self.handler_list[request.Method]
	response.Id = request.Id
	if nil == handler {
		response.Status = ERR_PROTOCOL
		response.Error = "undefined api:" + request.Method
		self.sendResponse(conn,response)
		return
	}

	fv = reflect.ValueOf(handler)

	call_params = request.Params.([]interface{})

	if len(call_params) != fv.Type().NumIn() {

		response.Status = ERR_REQUEST
		response.Error = "mismatch call param value type or length."
		self.sendResponse(conn,response)
		return

	}

	real_params = make([]reflect.Value, len(call_params))

	func() {

		for i, v := range call_params {
			raw_val := reflect.ValueOf(v)
			real_params[i] = raw_val.Convert(fv.Type().In(i))
		}

		rs = fv.Call(real_params)
	}()

	if len(rs) < 1 {

		response.Return(nil)

		self.sendResponse(conn,response)
		return
	}

	if len(rs) > 1 {

		response.Error = "Not Supported Multi Return Values"
		response.Status = ERR_OUTPUT
		self.sendResponse(conn,response)
		return
	}

	response.Return(rs[0].Interface())
	self.sendResponse(conn,response)

}

func (self *Server) Run() {
	defer self.transport.Stop()
	self.transport.OnConnection(self.handler)
	self.transport.Serve()
}
