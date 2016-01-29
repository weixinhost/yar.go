package yar

import (
	"bytes"
	"net"
	"yar/packager"
	"yar/transports"
	"fmt"
)

type Server struct {

	handler_list map[string]Handler

	opt map[string]interface{}

	transport transports.Transport
}

var CONNECTION_TOTAL = 0

func NewServer(host string,port int) *Server {

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

		size, err := self.transport.Read(conn, protocol_buffer[read_total:])

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

func (self *Server) handler(conn net.Conn) {

	var read_total int
	var body_buffer []byte
	var handler Handler = nil

	response := new(Response)
	request := new(Request)

	protocol, err := self.ReadProtocol(conn)

	if err != nil {
		response.Status = ERR_PROTOCOL
		response.Error = err.Error()
		goto send
	}

	body_buffer = make([]byte, protocol.BodyLength)

	for {

		len, err := conn.Read(body_buffer[read_total:])

		if err != nil {
			response.Status = ERR_PROTOCOL
			response.Error = err.Error()
			goto send
			break
		}

		read_total += len
		if read_total >= int(protocol.BodyLength) {
			break
		}
	}

	fmt.Printf("%s\n",body_buffer)
	fmt.Println(body_buffer)
	err = packager.Unpack(protocol.Packager[0:], body_buffer, &request)

	if err != nil {
		response.Status = ERR_PACKAGER
		response.Error = err.Error()
		goto send
	}

	request.Protocol = protocol
	request.Id = request.Protocol.Id
	handler = self.handler_list[request.Method]
	response.Id = request.Id
	response.Protocol = protocol

	if nil == handler {
		response.Status = ERR_PROTOCOL
		response.Error = "undefined api:" + request.Method
		goto send
	}

	handler(request,response)

send:

	ret, err := packager.Pack(protocol.Packager[0:], response)

	protocol.BodyLength = uint32(len(ret) + 8)
	send_protocol := protocol.Bytes()
	self.transport.Write(conn, send_protocol.Bytes())
	self.transport.Write(conn, ret)
	conn.Close()

}

func (self *Server) Run() {

	defer self.transport.Stop()
	self.transport.OnConnection(self.handler)
	self.transport.Run()

}
