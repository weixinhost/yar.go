package transports

import (
	"net"
	"os"
	"time"
)


type UdpConnection struct {
	conn net.Conn
}

func newUdpConnection(conn net.Conn)(*UdpConnection){
	tcpConn := new(UdpConnection)
	tcpConn.conn = conn
	return tcpConn
}

func (conn *UdpConnection)Read(buffer []byte)(n int ,err error){
	return conn.conn.Read(buffer)
}

func (conn *UdpConnection)Write(buffer[]byte)(n int ,err error){
	return conn.conn.Write(buffer)
}

func (conn *UdpConnection)Close()(err error){
	return conn.conn.Close()
}

func (conn *UdpConnection)SetReadTimeout(timeout time.Duration){
	now := time.Now()
	conn.conn.SetReadDeadline(now.Add(timeout))
}

func (conn *UdpConnection)SetWriteTimeout(timeout time.Duration){
	now := time.Now()
	conn.conn.SetWriteDeadline(now.Add(timeout))
}


type Udp struct {
	hostname string
	listener net.Listener
	handler ConnectionHandler
	running bool
}

func NewUdp(hostname string) (*Udp, error) {
	tcp := new(Udp)
	tcp.hostname = hostname
	tcp.handler = defaultHandler
	return tcp, nil
}

func (self *Udp) OnConnection(handler ConnectionHandler) {

	self.handler = handler
}

func (self *Udp) Serve()(err error) {

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

		tcpConn := newUdpConnection(conn)
		self.initConnection(tcpConn)
		go self.handler(tcpConn)
	}

	return nil
}

func(self *Udp)Connection()(t TransportConnection,err error){
	conn,err  :=  net.Dial("tcp",self.hostname)
	tcpConn := newTcpConnection(conn)
	self.initConnection(tcpConn)
	return tcpConn,err
}

func (self *Udp)initConnection(conn TransportConnection){
	conn.SetReadTimeout(CONNECTION_READ_TIMEOUT_SECOND * time.Second)
	conn.SetWriteTimeout(CONNECTION_READ_TIMEOUT_SECOND * time.Second)
}

func (self *Udp) Stop() {

	self.running = false

}
