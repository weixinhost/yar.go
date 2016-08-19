package client

import (
	"net"
	"sync"
	"time"
)

type ResolverResult struct {
	list    []net.IP
	Expired int
}

type Resolver struct {
	cache  map[string]ResolverResult
	lock   sync.RWMutex
	Expire time.Duration
	Max    int
}

func NewResolver(max int, expires time.Duration) *Resolver {
	r := new(Resolver)
	r.Max = max
	r.Expire = expires
	r.cache = make(map[string]ResolverResult)
	return r
}

func (r *Resolver) Lookup(domain string) ([]net.IP, error) {
	now := time.Now().Unix()
	r.lock.RLock()
	ret, ok := r.cache[domain]
	r.lock.RUnlock()
	if ok {
		if ret.Expired == 0 {
			return ret.list, nil
		}
		if ret.Expired > int(now) {
			return ret.list, nil
		}
	}
	ips, err := net.LookupIP(domain)
	if err != nil {
		return nil, err
	}
	ret = ResolverResult{
		list: ips,
	}
	if r.Expire > 0 {
		ret.Expired = int(now) + int(r.Expire.Seconds())
	}

	r.lock.Lock()
	r.cache[domain] = ret
	if r.Max > 0 && r.Max < len(r.cache) {
		i := 0
		newCache := make(map[string]ResolverResult)
		for k, item := range r.cache {
			newCache[k] = item
			if i > r.Max {
				break
			}
		}
		r.cache = newCache
	}
	r.lock.Unlock()
	return ips, err
}

var globalResolver *Resolver = nil

func init() {
	globalResolver = NewResolver(1000, 60*time.Second)
}
