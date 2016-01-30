package yar
import (
	"math/rand"
	"yar/packager"
	"net"
	"fmt"
	"bytes"
	"errors"
)

type Opt int

const (

	CONNECTION_TIMEOUT Opt = 1
	TIMEOUT               Opt = 2
	PACKAGER           Opt = 3

)

const (

	TCP_CLIENT = 1
	HTTP_CLIENT = 2
	UDP_CLIENT = 3

)

const (

	DEFAULT_PACKAGER = "json"
	DEFAULT_TIMEOUT = 5000
	DEFAULT_CONNECTION_TIMEOUT = 1000

)

type Client struct {

	netmode  int
	hostname string
	request  *Request
	opt      map[Opt]interface{}

}

func NewClientWithTcp(host string, port int) (client *Client, err error) {

	client = new(Client)
	client.request = new(Request)
	client.opt = make(map[Opt]interface{})
	client.request.Protocol = NewProtocol()
	client.hostname = fmt.Sprintf("%s:%d", host, port)
	client.netmode = TCP_CLIENT
	client.initOpt()

	return client, nil

}


func (self *Client)initOpt() {

	self.opt[CONNECTION_TIMEOUT] = DEFAULT_CONNECTION_TIMEOUT
	self.opt[TIMEOUT] = DEFAULT_TIMEOUT
	self.opt[PACKAGER] = DEFAULT_PACKAGER

}

func (self *Client)SetOpt(opt Opt, v interface{}) bool {

	switch opt {

	case CONNECTION_TIMEOUT:
	case TIMEOUT:
	case PACKAGER:{
		self.opt[opt] = v
		return true
	}

	}

	return false
}

func (self *Client)tcpCall(method string, ret interface{}, params ... interface{}) (err error) {

	if params != nil {
		self.request.Params = params
	}else {
		self.request.Params = []string{}
	}
	self.request.Id = rand.Uint32()
	self.request.Method = method
	self.request.Protocol.Id = self.request.Id
	self.request.Protocol.MagicNumber = MAGIC_NUMBER

	var pack []byte

	if len(self.opt[PACKAGER].(string)) < 8 {

		for i := 0; i < len(self.opt[PACKAGER].(string)); i++ {

			self.request.Protocol.Packager[i] = self.opt[PACKAGER].(string)[i]

		}

	}

	pack, err = packager.Pack([]byte(self.opt[PACKAGER].(string)), self.request)

	if err != nil {

		return err

	}

	self.request.Protocol.BodyLength = uint32(len(pack) + 8)


	conn, conn_err := net.Dial("tcp", self.hostname)

	if conn_err != nil {

		return conn_err
	}


	conn.Write(self.request.Protocol.Bytes().Bytes())
	conn.Write(pack)
	protocol_buffer := make([]byte, PROTOCOL_LENGTH)
	conn.Read(protocol_buffer)
	self.request.Protocol.Init(bytes.NewBuffer(protocol_buffer))
	body_buffer := make([]byte, self.request.Protocol.BodyLength - 8)
	conn.Read(body_buffer)
	response := new(Response)
	err = packager.Unpack([]byte(self.opt[PACKAGER].(string)), body_buffer, &response)

	if response.Status != ERR_OKEY {
		return errors.New(response.Error)
	}

	if ret != nil {

		err = packager.Unpack([]byte(self.opt[PACKAGER].(string)), bytes.NewBufferString(response.Retval).Bytes(), ret)
		return err
	}
	return nil
}

func (self *Client)Call(method string, ret interface{}, params ... interface{}) (err error) {

	switch self.netmode {

	case TCP_CLIENT : {
		return self.tcpCall(method, ret, params...)
	}
	}
	return errors.New("unsupported client netmode")
}
