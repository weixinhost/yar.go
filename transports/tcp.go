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

	listener net.Listener

	handler ConnectionHandler

	running bool
}

func defaultHandler(conn net.Conn) {

	conn.Close()
}

func NewTcp(host string, port int) (*Tcp, error) {

	tcp := new(Tcp)

	hostname := fmt.Sprintf("%s:%d",host,port)

	listener, err := net.Listen(NetMode, hostname)

	if err != nil {
		return nil, err
	}

	tcp.listener = listener

	tcp.handler = defaultHandler

	tcp.running = true

	return tcp, nil
}

func (self *Tcp) OnConnection(handler ConnectionHandler) {

	self.handler = handler
}

func (self *Tcp) Run() {

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

}

func (self *Tcp) Read(conn net.Conn, buffer []byte) (len int, err error) {

	len, err = conn.Read(buffer)

	return len, err

}

func (self *Tcp) Write(conn net.Conn, buffer []byte) (real_len int, err error) {

	real_len, err = conn.Write(buffer)

	return real_len, err
}

func (self *Tcp) Stop() {

	self.running = false

}
