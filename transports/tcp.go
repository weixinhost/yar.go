package transports

import (
	"net"
	"os"
	"time"
)


type TcpConnection struct {
	conn net.Conn
}

func newTcpConnection(conn net.Conn)(*TcpConnection){
	tcpConn := new(TcpConnection)
	tcpConn.conn = conn
	return tcpConn
}

func (conn *TcpConnection)Read(buffer []byte)(n int ,err error){
	return conn.conn.Read(buffer)
}

func (conn *TcpConnection)Write(buffer[]byte)(n int ,err error){
	return conn.conn.Write(buffer)
}

func (conn *TcpConnection)Close()(err error){
	return conn.conn.Close()
}

func (conn *TcpConnection)SetReadTimeout(timeout time.Duration){
		now := time.Now()
		conn.conn.SetReadDeadline(now.Add(timeout))
}

func (conn *TcpConnection)SetWriteTimeout(timeout time.Duration){
	now := time.Now()
	conn.conn.SetWriteDeadline(now.Add(timeout))
}

type Tcp struct {
	hostname string
	listener net.Listener
	handler ConnectionHandler
	running bool
}

func NewTcp(hostname string) (*Tcp, error) {
	tcp := new(Tcp)
	tcp.hostname = hostname
	tcp.handler = defaultHandler
	return tcp, nil
}

func (self *Tcp) OnConnection(handler ConnectionHandler) {

	self.handler = handler
}



func (self *Tcp) Serve()(err error) {

	listener, err := net.Listen("tcp", self.hostname)

	if err != nil {
		return err
	}

	self.listener = listener
	self.running = true

	defer self.listener.Close()

	for {

		if self.running == false {
			break
		}

		conn, err := self.listener.Accept()

		if err != nil {
			os.Exit(-1)
		}

		tcpConn := newTcpConnection(conn)
		self.initConnection(tcpConn)
		go self.handler(tcpConn)
	}

	return nil

}

func(self *Tcp)Connection()(t TransportConnection,err error){
	conn,err  :=  net.Dial("tcp",self.hostname)
	tcpConn := newTcpConnection(conn)
	self.initConnection(tcpConn)
	return tcpConn,err
}

func (self *Tcp)initConnection(conn TransportConnection){

	conn.SetReadTimeout(CONNECTION_READ_TIMEOUT_SECOND * time.Second)
	conn.SetWriteTimeout(CONNECTION_READ_TIMEOUT_SECOND * time.Second)

}

func (self *Tcp) Stop() {

	self.running = false

}
