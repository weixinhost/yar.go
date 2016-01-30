package transports

import (
	"net"
	"os"
	"fmt"
)

const (
	NetMode = "tcp"
)

type Tcp struct {

	host string
	port int
	listener net.Listener
	handler ConnectionHandler
	running bool
}

func defaultHandler(conn net.Conn) {

	conn.Close()
}

func NewTcp(host string, port int) (*Tcp, error) {
	tcp := new(Tcp)
	tcp.host = host
	tcp.port = port
	tcp.handler = defaultHandler
	return tcp, nil
}

func (self *Tcp) OnConnection(handler ConnectionHandler) {

	self.handler = handler
}

func (self *Tcp) Serve()(err error) {

	hostname := fmt.Sprintf("%s:%d",self.host,self.port)
	listener, err := net.Listen(NetMode, hostname)

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
		go self.handler(conn)
	}

	return nil

}

func(self *Tcp)Connection()(conn net.Conn,err error){
	hostname := fmt.Sprintf("%s:%d",self.host,self.port)
	return net.Dial("tcp",hostname)
}

func (self *Tcp) Stop() {

	self.running = false

}
