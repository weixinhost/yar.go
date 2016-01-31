package yar
import (
	"strings"
	"yar/transports"
	"net"
	"bytes"
	"errors"
	"yar/packager"
	"reflect"
	"fmt"
)


type Server struct {

	opt 		map[string]interface{}
	handlerList map[string]Handler
	netProtocol string
	hostname 	string
	transport 	transports.Transport
}

func NewServer(net string,hostname string) (server *Server,err error) {

	server = new(Server)
	server.handlerList = make(map[string]Handler,32)
	server.netProtocol = net
	server.hostname = hostname
	server.transport = nil

	return
}


func (server *Server) RegisterHandler(name string ,handler Handler){

	server.handlerList[name] = handler

}

func (server *Server) RemoveHandler(name string) {

	delete(server.handlerList,name)

}

func (server *Server) OnConnection(conn net.Conn) {

	protocol,err := server.getProtocol(conn)
	response := NewResponse()

	if err != nil {
		server.sendResponse(conn,response)
		return ;
	}

	response.Protocol = protocol

	request,err := server.getRequest(conn,protocol)
	if err != nil {
		response.Error = "get or parse request errror"
		response.Status = ERR_TRANSPORT
		server.sendResponse(conn,response)
		return
	}

	request.Id 	= protocol.Id
	response.Id = request.Id

	request.Protocol = protocol
	server.call(request,response)
	server.sendResponse(conn,response)
}


func (server *Server) Serve() {

	server.init()

	if server.transport == nil {
		panic("server.transport is nil")
	}

	server.transport.OnConnection(server.OnConnection)

	server.transport.Serve()
}



// ===================== private =======================


func (server *Server)init(){

	switch strings.ToLower(server.netProtocol) {

	case "tcp" 	: {

		server.transport,_ = transports.NewTcp(server.hostname)

		break
	}

	case "udp" 	:{
		break
	}

	case "http" :{
		break
	}

	case "unix" 	:{

		break
	}

	}

}


func (server *Server)getProtocol(conn net.Conn)(protocol *Protocol,err error){

	protocolBuffer := make([]byte,PROTOCOL_LENGTH + PACKAGER_LENGTH)

	receiveLen := 0

	for {

		realLen,err := conn.Read(protocolBuffer)

		if err != nil {
			return nil,err
		}

		receiveLen += realLen

		if receiveLen >= len(protocolBuffer) {
			break
		}

	}

	protocol = NewProtocolWithBytes(bytes.NewBuffer(protocolBuffer))

	return protocol,nil

}


func (server *Server)getRequest(conn net.Conn,protocol *Protocol)(request *Request,err error){

	realBodyLen := protocol.BodyLength - PACKAGER_LENGTH

	if realBodyLen < 0 {
		return nil,errors.New("protocol body length parse error")
	}

	bodyBuffer := make([]byte,realBodyLen)
	receiveLen := 0

	for {

		realLen,err := conn.Read(bodyBuffer)

		if err != nil {
			return nil,err
		}

		receiveLen += realLen

		if receiveLen >= int(realBodyLen) {
			break
		}
	}

	request = NewRequest()
	err = packager.Unpack(protocol.Packager[:],bodyBuffer,request)
	return request,err
}


func (server *Server)sendResponse(conn net.Conn,response *Response)(err error){

	defer conn.Close()

	if response.Protocol != nil {

		sendPackData,err := packager.Pack(response.Protocol.Packager[:],response)

		if err != nil {
				return err
		}

		response.Protocol.BodyLength = uint32(len(sendPackData) +8)

		conn.Write(response.Protocol.Bytes().Bytes())
		conn.Write(sendPackData)
	}

	return nil
}


func (server *Server)call(request *Request,response *Response) {

	defer func(){

		if err := recover(); err != nil {
			response.Status = ERR_EMPTY_RESPONSE
			response.Error = fmt.Sprintf("server has panic:%s",err)
		}
	}()

	handler := server.handlerList[request.Method]

	if handler == nil {
		response.Status = ERR_EMPTY_RESPONSE
		response.Error = "call undefined api:" + request.Method
		return
	}
	call_params := request.Params.([]interface{})

	fv := reflect.ValueOf(handler)

	if len(call_params) != fv.Type().NumIn() {

		response.Status = ERR_EMPTY_RESPONSE
		response.Error = "mismatch call params"
		return
	}

	real_params := make([]reflect.Value, len(call_params))

	func() {

		for i, v := range call_params {
			raw_val := reflect.ValueOf(v)
			real_params[i] = raw_val.Convert(fv.Type().In(i))
		}

		rs := fv.Call(real_params)
		if len(rs) < 1 {
			response.Return(nil)
		}

		if len(rs) > 1 {
			response.Status = ERR_EMPTY_RESPONSE
			response.Error = "unsupprted multi value return on rpc call"
			return
		}

		response.Return(rs[0].Interface())
	}()
}

