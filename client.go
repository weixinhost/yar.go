package yar
import (
	"yar/transports"
	"math/rand"
)

type Opt int

const (

	CONNECTION_TIMEOUT = 1
	TIMEOUT			   = 2
	PACKAGER 		   = 3

)

type Client struct {

	request *Request
	opt 	map[int]interface{}
	transports transports.Transport

}

func NewClientWithTcp(host string,port int) (client *Client,err error) {

	client = new(Client)
	client.request = new(Request)
	client.opt = make(map[int]interface{})
	client.transports = transports.NewTcp(host,port)
	return client,nil

}

func (self *Client)SetOpt(opt Opt, v interface{}) bool {

	switch opt {

	case CONNECTION_TIMEOUT || TIMEOUT || PACKAGER:{
		self.opt[opt] = v
		return true
	}

	}

	return false
}

func (self *Client)Call(method string, params ... interface{}) (ret interface{}) {

	self.request.Id = rand.Uint32()

	self.request.Method = method

	self.request.Params = params

	return nil
}


