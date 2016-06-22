package transports

import (
	"net"
	"os"
	"time"
)

type SockConnection struct {
	conn net.Conn
}

func newSockConnection(conn net.Conn) *SockConnection {
	tcpConn := new(SockConnection)
	tcpConn.conn = conn
	return tcpConn
}

func (conn *SockConnection) Read(buffer []byte) (n int, err error) {
	return conn.conn.Read(buffer)
}

func (conn *SockConnection) Write(buffer []byte) (n int, err error) {
	return conn.conn.Write(buffer)
}

func (conn *SockConnection) Close() (err error) {
	return conn.conn.Close()
}

func (conn *SockConnection) SetReadTimeout(timeout time.Duration) {
	now := time.Now()
	conn.conn.SetReadDeadline(now.Add(timeout))
}

func (conn *SockConnection) SetWriteTimeout(timeout time.Duration) {
	now := time.Now()
	conn.conn.SetWriteDeadline(now.Add(timeout))
}

type Sock struct {
	hostname string
	net      string
	listener net.Listener
	handler  ConnectionHandler
	running  bool
}

func NewSock(net string, hostname string) (*Sock, error) {
	tcp := new(Sock)
	tcp.hostname = hostname
	tcp.handler = defaultHandler
	tcp.net = net
	return tcp, nil
}

func (self *Sock) OnConnection(handler ConnectionHandler) {

	self.handler = handler
}

func (self *Sock) Serve() (err error) {

	listener, err := net.Listen(self.net, self.hostname)

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

		tcpConn := newSockConnection(conn)
		self.initConnection(tcpConn)
		go self.handler(tcpConn)
	}

	return nil

}

func (self *Sock) Connection() (t TransportConnection, err error) {
	conn, err := net.Dial(self.net, self.hostname)

	if err != nil {

		return nil, err
	}

	tcpConn := newSockConnection(conn)
	self.initConnection(tcpConn)
	return tcpConn, err
}

func (self *Sock) initConnection(conn TransportConnection) {

	conn.SetReadTimeout(CONNECTION_READ_TIMEOUT_SECOND * time.Second)
	conn.SetWriteTimeout(CONNECTION_READ_TIMEOUT_SECOND * time.Second)

}

func (self *Sock) Stop() {

	self.running = false

}
