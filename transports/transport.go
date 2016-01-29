package transports

import "net"

type ConnectionHandler func(conn net.Conn)

type Transport interface {
	Run()
	Stop()
	OnConnection(handler ConnectionHandler)
	Read(conn net.Conn, buffer []byte) (len int, err error)
	Write(conn net.Conn, buffer []byte) (real_len int, err error)
}
