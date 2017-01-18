package distribunted_client

import (
	"fmt"

	"strings"

	yar "github.com/weixinhost/yar.go"
	"github.com/weixinhost/yar.go/client"
	"github.com/weixinhost/yar.go/host_sync"
)

var pool *PeerPool

type Client struct {
	pool     string
	name     string
	path     string
	protocol string
}

func NewClient(protocol, pool, name, path string) *Client {
	c := new(Client)
	c.protocol = protocol
	c.name = name
	c.pool = pool
	c.path = path
	return c
}

func (self *Client) Call(method string, ret interface{}, params ...interface{}) *yar.Error {
	p := pool.GetPeer(self.pool, self.name)
	c := 0
	for {
		c++
		host, err := p.GetNextHost()
		if p.Len() <= c {
			break
		}
		if err != nil {
			return yar.NewError(yar.ErrorNetwork, "GetHostFailed:"+err.Error())
		}
		u := fmt.Sprintf("%s://%s/%s", self.protocol, host, self.path)
		c, aerr := client.NewClient(u)
		if aerr != nil {
			return aerr
		}
		e := c.Call(method, ret, params...)
		if e == nil {
			p.Reset(host)
			return nil
		}

		if e.Assert(yar.ErrorNetwork) && strings.Contains(e.String(), "connection refused") {
			p.SetFail(host)
		} else {
			p.Reset(host)
			return e
		}
	}
	return yar.NewError(yar.ErrorNetwork, "No Health Service Found.")
}

func Setup(dockerAPI, redisHost string) {
	host_sync.SetDockerAPI(dockerAPI)
	host_sync.SetRedisHost(redisHost)
}

func init() {
	pool = NewPeerPool()
}
