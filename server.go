package yar

import (
	"yar/transports"
	"net"
	"bytes"
//	"fmt"
//	"fmt"
	//"strings"
	//"yar/packager"
	"yar/packager"
	"fmt"
)

type Server struct {

	handler_list map[string]Handler

	transport transports.Transport

}

var CONNECTION_TOTAL = 0

func NewServer() *Server {

	s := new(Server)

	return s

}


func (self *Server)RegisterHandler(name string,handler Handler) bool {

	self.handler_list[name] = handler

	return true
}

func (self *Server)RemoveHandler(name string) bool {

	delete(self.handler_list,name)

	return true
}

func (self *Server)ReadProtocol(conn net.Conn) (*Protocol,error){

	protocol := NewProtocol()

	protocol_len := PROTOCOL_LENGTH

	protocol_buffer := make([]byte,protocol_len)

	read_total := 0

	for {

		size,err := self.transport.Read(conn,protocol_buffer[read_total:])

		if err != nil {

			return nil,err
		}

		read_total += size

		if read_total >= int(protocol_len) {

			break

		}

	}

	buffer := bytes.NewBuffer(protocol_buffer)

	protocol.Init(buffer)

	protocol.BodyLength -= 8

	return protocol,nil
}

func (self *Server)handler(conn net.Conn) {

	response := new(Response)

	protocol,err := self.ReadProtocol(conn)

	if err != nil {

		//todo error handler

	}

	read_total := 0

	body_buffer := make([]byte,protocol.BodyLength)

	for {

		len,err := conn.Read(body_buffer[read_total:])

		if err != nil {

			//todo error handler

			break;

		}

		read_total += len

		if read_total >= int(protocol.BodyLength) {

			break;

		}
	}


	packager_buffer := bytes.NewBufferString("")

	for i:= 0;i<len(protocol.Packager);i++ {

		if protocol.Packager[i] == byte(0) {

			break;

		}

		packager_buffer.WriteByte(protocol.Packager[i])

	}

	request := new(Request)

	err = packager.Unpack(packager_buffer.String(),body_buffer,&request)

	if err != nil {

		response.Status = ERR_PACKAGER
		response.Error = err.Error()
		goto send
	}

	request.Protocol = protocol

	handler := self.handler_list[request.Method]

	response.Id = request.Id
	response.Protocol = protocol

	if nil == handler {

		response.Status = ERR_PROTOCOL
		response.Error = "undefined api:" + request.Method
		response.Output = response.Error

		goto send

	}

	send :

	ret,_ := packager.Pack(packager_buffer.String(),&response)

	response.Protocol.BodyLength = uint32(len(ret) + 8)

	send_protocol := response.Protocol.Bytes()


	self.transport.Write(conn,send_protocol.Bytes())

	self.transport.Write(conn,ret)


	 conn.Close()



}


func (self *Server)Run(){

	tran  ,_ := transports.NewTcp("0.0.0.0","6789")

	self.transport = tran

	self.transport.OnConnection(self.handler)

	self.transport.Run()

}
