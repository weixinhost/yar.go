package transports

import "net"

const (

	CONNECTION_READ_TIMEOUT_SECOND 	= 5
	CONNECTION_WRITE_TIMEOUT_SECOND = 5
)



type ConnectionHandler func(conn net.Conn)


type Transport interface {
	Serve()(err error)
	OnConnection(handler ConnectionHandler)
	Connection()(conn net.Conn,err error)
}

