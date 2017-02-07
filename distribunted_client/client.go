package distribunted_client

import (
	"fmt"
	"time"

	"strings"

	yar "github.com/weixinhost/yar.go"
	"github.com/weixinhost/yar.go/client"
	"github.com/weixinhost/yar.go/host_sync"
	"github.com/weixinhost/yar.go/monitor"
)

var pool *PeerPool

type Client struct {
	pool     string
	name     string
	path     string
	protocol string
	Opt      *yar.Opt
}

func NewClient(protocol, pool, name, path string) *Client {
	c := new(Client)
	c.protocol = protocol
	c.name = name
	c.pool = pool
	c.path = path
	c.Opt = yar.NewOpt()
	return c
}

func (self *Client) Call(method string, ret interface{}, params ...interface{}) *yar.Error {
	p := pool.GetPeer(self.pool, self.name)
	c := 0
	for {
		c++
		host, err := p.GetNextHost()

		if err != nil {
			return yar.NewError(yar.ErrorNetwork, err.Error())
		}

		if p.Len() <= c {
			break
		}

		u := fmt.Sprintf("%s://%s/%s", self.protocol, host, self.path)
		c, aerr := client.NewClient(u)
		if aerr != nil {
			return aerr
		}
		var opt yar.Opt
		opt = *self.Opt
		c.Opt = &opt
		now := time.Now()
		e := c.Call(method, ret, params...)
		end := time.Now()
		mils := int(end.Sub(now).Seconds() * 1000)

		if e == nil {
			monitor.SetServiceMonitor(self.pool, self.name, opt.Provider, mils, true)
			p.Reset(host)
			return nil
		}

		monitor.SetServiceMonitor(self.pool, self.name, opt.Provider, mils, false)

		//mismatch service name or host down
		if e.Assert(yar.ErrorNetwork) && strings.Contains(e.String(), "connection refused") || (e.Assert(yar.ErrorProtocol) && strings.Contains(e.String(), "mismatch service name")) {
			p.SetFail(host)
		} else {
			p.Reset(host)
			return e
		}
	}
	return yar.NewError(yar.ErrorNetwork, "No Health Service Found.")
}

func Setup(dockerAPI, redisHost string, h monitor.RealTimeMonitorHandle) {
	host_sync.SetDockerAPI(dockerAPI)
	host_sync.SetRedisHost(redisHost)
	monitor.Setup(redisHost, h)
}

func init() {
	pool = NewPeerPool()
}
