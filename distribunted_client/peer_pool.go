package distribunted_client

import (
	"fmt"
	"sync"
)

type PeerPool struct {
	pool  map[string]*Peer
	mutex sync.Mutex
}

func NewPeerPool() *PeerPool {
	pool := new(PeerPool)
	pool.pool = make(map[string]*Peer)
	return pool
}

func (p *PeerPool) GetPeer(pool, name string) *Peer {

	n := fmt.Sprintf("%s_%s", pool, name)

	p.mutex.Lock()
	peer, ok := p.pool[n]
	if !ok {
		peer = NewPeer(pool, name)
		p.pool[n] = peer
	}
	p.mutex.Unlock()
	return peer
}

func (p *PeerPool) Close() {}
