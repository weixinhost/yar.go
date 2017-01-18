package client

import (
	"errors"
	"net"
	"sync"
	"time"
)

var PeerErrorNotIpFound error = errors.New("peer not found alive ip")

const (
	ConnStatusIdle int = 0
	ConnStatusWork int = 1
)

type Conn struct {
	addr     string
	conn     *net.TCPConn
	count    int
	connTime time.Time
	lastTime time.Time
}

func NewConn(addr string) (*Conn, error) {
	conn := new(Conn)
	conn.addr = addr
	c, err := net.DialTimeout("tcp", addr, 5*time.Second)
	if err != nil {
		return nil, err
	}
	conn.conn = c.(*net.TCPConn)
	conn.conn.SetKeepAlive(true)
	conn.conn.SetKeepAlivePeriod(60 * time.Second)
	conn.connTime = time.Now()
	return conn, nil
}

type Peer struct {
	name     string
	ips      []string
	conns    map[string][]*Conn
	mutex    sync.Mutex
	maxIdles int
	offset   int
}

func NewPeer(name string) *Peer {
	peer := new(Peer)
	peer.name = name
	peer.ips = make([]string, 0)
	return peer
}

func (p *Peer) SyncIpList(ips []string) {
	p.mutex.Lock()
	p.ips = ips
	p.mutex.Unlock()
}

func (p *Peer) GetConn() (*Conn, error) {
	p.mutex.Lock()
	if p.ips == nil || len(p.ips) < 1 {
		p.mutex.Unlock()
		return nil, PeerErrorNotIpFound
	}
	if p.offset >= len(p.ips) {
		p.offset = 0
	}
	addr := p.ips[p.offset]
	p.offset++
	if _, ok := p.conns[addr]; ok {
		if len(p.conns[addr]) > 1 {
			l := len(p.conns[addr])
			conn := p.conns[addr][l-1]
			p.conns[addr] = p.conns[addr][:l-1]
			p.mutex.Unlock()
			return conn, nil
		}
	}
	p.mutex.Unlock()
	conn, err := NewConn(addr)
	if err != nil {
		return nil, err
	}
	p.PushConn(conn)
	return conn, err
}

func (p *Peer) PushConn(conn *Conn) error {
	addr := conn.conn.RemoteAddr().String()
	p.mutex.Lock()
	if _, ok := p.conns[addr]; !ok {
		p.conns[addr] = make([]*Conn, 0)
	}
	p.conns[addr] = append(p.conns[addr], conn)
	p.mutex.Unlock()
	return nil
}
