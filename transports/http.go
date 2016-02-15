package transports

import (
	"net"
	//"os"
	"time"
	"net/http"
	"fmt"
)

type HttpConnection struct {
	request *http.Request
	response http.ResponseWriter
	body []byte
	use_read int
	requestTime time.Time
	responseTime time.Time
}


func NewHttpConnection(response http.ResponseWriter,request *http.Request)(*HttpConnection){

	conn := new(HttpConnection)
	conn.request = request
	conn.response = response
	conn.use_read = 0

	return conn
}

func(conn *HttpConnection)Read(buffer []byte)(n int,err error){

	return conn.request.Body.Read(buffer)
}

func(conn *HttpConnection)Write(buffer []byte)(n int,err error){
	return conn.response.Write(buffer)
}

func(conn *HttpConnection)Close()(err error){
//	 conn.request.Body.Close()
	 return nil
}

func(conn *HttpConnection)SetReadTimeout(timeout time.Duration){
	//its empty.
}

func (conn *HttpConnection)SetWriteTimeout(timeout time.Duration){
	//its empty.
}


func (conn *HttpConnection)GetRemoteAddr() string{
	addr := conn.request.RemoteAddr
	return addr
}


func (conn *HttpConnection)SetRequestTime(t time.Time){
	conn.requestTime = t
}


func (conn *HttpConnection)SetResponseTime(t time.Time){
	conn.responseTime = t
}


func (conn *HttpConnection)GetRequestTime()(t time.Time){
	return conn.requestTime
}


func (conn *HttpConnection)GetResponseTime()(t time.Time){
	return conn.responseTime
}



type Http struct {
	hostname string
	path string
	listener net.Listener
	handler ConnectionHandler
	readTimeout time.Duration
	writeTimeout time.Duration
	running bool
}

func NewHttp(hostname string,path string,readTimeout time.Duration,writeTimeout time.Duration) (*Http, error) {
	tcp := new(Http)
	tcp.hostname = hostname
	tcp.handler = defaultHandler
	tcp.path = path
	tcp.writeTimeout = writeTimeout
	tcp.readTimeout = readTimeout
	return tcp, nil
}

func (self *Http) OnConnection(handler ConnectionHandler) {

	self.handler = handler
}

func (self *Http)ServeHTTP(writer http.ResponseWriter,request *http.Request) {
	conn := NewHttpConnection(writer,request)
	self.handler(conn)
}

func (self *Http) Serve()(err error) {

	s := &http.Server{
		Addr : self.hostname,
		ReadTimeout:self.readTimeout,
		WriteTimeout:self.writeTimeout,
		Handler:self,
	}

	err = s.ListenAndServe()
	fmt.Print(err)
	return nil
}


func(self *Http)Connection()(conn TransportConnection,err error){

	return nil,err
}


func (self *Http)initConnection(conn net.Conn){

	now := time.Now()

	readDeadline := now.Add(CONNECTION_READ_TIMEOUT_SECOND * time.Second)
	conn.SetReadDeadline(readDeadline)

	writeDeadline := readDeadline.Add(CONNECTION_READ_TIMEOUT_SECOND * time.Second)
	conn.SetWriteDeadline(writeDeadline)
}

func (self *Http) Stop() {

	self.running = false

}
