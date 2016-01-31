package transports

import (
	"net"
	"os"
	"time"
)

const (
	NetMode = "tcp"
)

type Tcp struct {

	hostname string
	listener net.Listener
	handler ConnectionHandler
	running bool
}

func defaultHandler(conn net.Conn) {

	conn.Close()
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

	listener, err := net.Listen(NetMode, self.hostname)

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

		self.initConnection(conn)
		go self.handler(conn)
	}

	return nil

}

func(self *Tcp)Connection()(conn net.Conn,err error){

	conn,err =  net.Dial("tcp",self.hostname)
	self.initConnection(conn)

	return conn,err

}

func (self *Tcp)initConnection(conn net.Conn){

	now := time.Now()

	readDeadline := now.Add(CONNECTION_READ_TIMEOUT_SECOND * time.Second)
	conn.SetReadDeadline(readDeadline)

	writeDeadline := readDeadline.Add(CONNECTION_READ_TIMEOUT_SECOND * time.Second)
	conn.SetWriteDeadline(writeDeadline)
}

func (self *Tcp) Stop() {

	self.running = false

}
