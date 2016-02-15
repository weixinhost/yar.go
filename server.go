package yar
import (
	"strings"
	"yar/transports"
//	"net"
	"bytes"
	"errors"
	"yar/packager"
	"reflect"
	"fmt"
	"time"
	"yar/log"
)

const (
	SERVER_OPT_LOG_PATH = "log"
	SERVER_OPT_READ_TIMEOUT = "read_timeout"
	SERVER_OPT_WRITE_TIMEOUT = "write_timeout"
)

// Yar 服务端
// [file://examples/snowflake_server.go]
type Server struct {

	opt 		map[string]interface{}
	handlerList map[string]Handler
	netProtocol string
	hostname 	string
	transport 	transports.Transport
	log 		log.Log
}

//NewServer is easy to create a yar server
func NewServer(net string,hostname string) (server *Server,err error) {

	server = new(Server)
	server.handlerList = make(map[string]Handler,32)
	server.netProtocol = net
	server.hostname = hostname
	server.transport = nil
	server.opt = make(map[string]interface{},32)
	server.initServer()
	return
}

func (server *Server)initServer(){

	server.SetOpt(SERVER_OPT_READ_TIMEOUT,time.Duration(5))
	server.SetOpt(SERVER_OPT_WRITE_TIMEOUT,time.Duration(5))
}

func (server *Server) RegisterHandler(name string ,handler Handler){

	server.handlerList[name] = handler

}

func(server *Server)SetOpt(name string,v interface{}){
	server.opt[name] = v
}

func(server *Server)GetOpt(name string) interface{} {

	return server.opt[name]
}

func (server *Server) RemoveHandler(name string) {

	delete(server.handlerList,name)

}

func (server *Server) OnConnection(conn transports.TransportConnection) {

	conn.SetRequestTime(time.Now())

	defer conn.Close()

	conn.SetReadTimeout(server.GetOpt(SERVER_OPT_READ_TIMEOUT).(time.Duration) * time.Second)
	conn.SetWriteTimeout(server.GetOpt(SERVER_OPT_WRITE_TIMEOUT).(time.Duration) * time.Second)

	protocol,err := server.getProtocol(conn)
	response := NewResponse()

	if err != nil {
		server.sendResponse(conn,nil,response)
		return ;
	}
	response.Protocol = protocol
	request,err := server.getRequest(conn,protocol)
	if err != nil {
		response.Error = "get or parse request errror:" + err.Error()
		response.Status = ERR_TRANSPORT
		server.sendResponse(conn,nil,response)
		return
	}

	request.Id 	= protocol.Id
	response.Id = request.Id
	request.Protocol = protocol

	server.call(request,response)
	server.sendResponse(conn,request,response)
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

	case "tcp" , "udp" , "unix": {

		server.transport,_ = transports.NewSock(server.netProtocol,server.hostname)
		break
	}

	case "http" :{
		server.transport,_ = transports.NewHttp(server.hostname,"/",
			server.GetOpt(SERVER_OPT_READ_TIMEOUT).(time.Duration) * time.Second,
			server.GetOpt(SERVER_OPT_WRITE_TIMEOUT).(time.Duration) * time.Second)
		break
	}

	}

	if server.opt[SERVER_OPT_LOG_PATH] != nil  {
		server.log ,_ = log.NewFileLog(server.opt[SERVER_OPT_LOG_PATH].(string))
	}
}


func (server *Server)getProtocol(conn transports.TransportConnection)(protocol *Protocol,err error){

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


func (server *Server)getRequest(conn transports.TransportConnection,protocol *Protocol)(request *Request,err error){

	realBodyLen := protocol.BodyLength - PACKAGER_LENGTH
	realBodyLen = 34
	if realBodyLen < 0 {
		return nil,errors.New("protocol body length parse error")
	}

	bodyBuffer := make([]byte,realBodyLen)
	receiveLen := 0
	for {
		realLen,err := conn.Read(bodyBuffer)
		if err != nil  && err.Error() != "EOF"{
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


func(server *Server)writeLog(conn transports.TransportConnection,level log.LogLevel,fmt string,params...interface{}){

	if server.log != nil {

	//	server.log.Append(conn,level,fmt,params...)
	}

}

func (server *Server)sendResponse(conn transports.TransportConnection, request *Request, response *Response) (err error) {

	conn.SetResponseTime(time.Now())

	if response.Protocol != nil {

		sendPackData, err := packager.Pack(response.Protocol.Packager[:], response)

		if err != nil {
			server.writeLog(conn, log.LOG_ERROR, "pack data error:%s", err.Error())
			return err
		}

		if response.Status != ERR_OKEY {
			if request != nil {
				server.writeLog(conn, log.LOG_ERROR, "request error:%d %s %s", request.Id, request.Method, response.Error)
			}else {
				server.writeLog(conn, log.LOG_ERROR, "request error:%s", response.Error)
			}
			return nil
		}

		response.Protocol.BodyLength = uint32(len(sendPackData) + 8)

		_, err = conn.Write(response.Protocol.Bytes().Bytes())
		_, err = conn.Write(sendPackData)

		if err != nil {
			server.writeLog(conn, log.LOG_ERROR, "response error:%d %s %s", request.Id, request.Method, err.Error())
		}else {
			server.writeLog(conn, log.LOG_NORMAL, "%d %s OKEY", request.Id, request.Method)
		}
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