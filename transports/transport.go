package transports

import "net"

type ConnectionHandler func(conn net.Conn)

type Transport interface {
	Serve()(err error)
	Stop()
	OnConnection(handler ConnectionHandler)
	Connection()(conn net.Conn,err error)
}

