package transports

import (

	"io"
	"time"
)

const (

	CONNECTION_READ_TIMEOUT_SECOND 	= 5
	CONNECTION_WRITE_TIMEOUT_SECOND = 5

)

type TransportConnection interface{
	io.Reader
	io.Writer
	io.Closer

	SetReadTimeout(timeout time.Duration)
	SetWriteTimeout(timeout time.Duration)
	GetRemoteAddr() 	string
	GetRequestTime() 	time.Time
	GetResponseTime() 	time.Time
	SetRequestTime(time time.Time)
	SetResponseTime(time time.Time)
}

type ConnectionHandler func(conn TransportConnection)

type Transport interface {
	Serve()(err error)
	OnConnection(handler ConnectionHandler)
	Connection()(conn TransportConnection,err error)
}

func defaultHandler(conn TransportConnection) {
	conn.Close()
}