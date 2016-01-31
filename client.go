package yar

import (
	"bytes"
	"errors"
	"math/rand"
	"yar/packager"
	"yar/transports"
	"encoding/gob"
	"strings"
	"fmt"
)

type Opt int

const (
	CONNECTION_TIMEOUT Opt = 1
	TIMEOUT            Opt = 2
	PACKAGER           Opt = 3
)

const (
	TCP_CLIENT  = 1
	HTTP_CLIENT = 2
	UDP_CLIENT  = 3
)

const (
	DEFAULT_PACKAGER           = "json"
	DEFAULT_TIMEOUT            = 5000
	DEFAULT_CONNECTION_TIMEOUT = 1000
)

type Client struct {
	net string
	hostname string
	request   *Request
	transport transports.Transport
	opt       map[Opt]interface{}
}

func NewClient(net string,hostname string)(client *Client){

	client = new(Client)
	client.hostname = hostname
	client.net = strings.ToLower(net)
	client.opt = make(map[Opt]interface{},6)
	client.request = NewRequest()
	client.request.Protocol = NewProtocol()
	client.initOpt()
	client.init()

	return client
}

func (client *Client)init() {

	switch client.net {

	case "tcp" : {
		client.transport,_ = transports.NewTcp(client.hostname)
		break
	}

	}

}

func (self *Client) initOpt() {

	self.opt[CONNECTION_TIMEOUT] = DEFAULT_CONNECTION_TIMEOUT
	self.opt[TIMEOUT] = DEFAULT_TIMEOUT
	self.opt[PACKAGER] = DEFAULT_PACKAGER

}

func (self *Client) SetOpt(opt Opt, v interface{}) bool {

	switch opt {

	case CONNECTION_TIMEOUT:
	case TIMEOUT:
	case PACKAGER:
		{
			self.opt[opt] = v
			return true
		}

	}

	return false
}

func (self *Client) tcpCall(method string,ret interface{},params ...interface{}) (err error) {

	if params != nil {
		self.request.Params = params
	} else {
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

	self.request.Protocol.BodyLength = uint32(len(pack) + PACKAGER_LENGTH)
	conn, conn_err := self.transport.Connection()

	if conn_err != nil {
		return conn_err
	}

	conn.Write(self.request.Protocol.Bytes().Bytes())
	conn.Write(pack)
	protocol_buffer := make([]byte, PROTOCOL_LENGTH + PACKAGER_LENGTH)
	conn.Read(protocol_buffer)
	self.request.Protocol.Init(bytes.NewBuffer(protocol_buffer))
	body_buffer := make([]byte, self.request.Protocol.BodyLength - PACKAGER_LENGTH)
	conn.Read(body_buffer)
	response := new(Response)
	err = packager.Unpack([]byte(self.opt[PACKAGER].(string)), body_buffer, &response)

	if response.Status != ERR_OKEY {
		return errors.New(response.Error)
	}

	fmt.Printf("%s\n",body_buffer)

	//这里需要优化,需要干掉这次pack/unpack
	pack_data,err := packager.Pack(self.request.Protocol.Packager[:],response.Retval)
	err = packager.Unpack(self.request.Protocol.Packager[:],pack_data,ret)

	return err
}

func (self *Client) Call(method string, ret interface{},params ...interface{}) (err error) {

	switch self.net {

	case "tcp":
		{
			return self.tcpCall(method, ret,params...)
		}
	}

	return errors.New("unsupported client netmode")
}

func (self *Client)parseRetVal(retval interface{},parse interface{})(err error){

	buf := bytes.NewBufferString("")

	enc := gob.NewEncoder(buf)
	dec := gob.NewDecoder(buf)

	err = enc.Encode(retval)

	if err != nil {

		return err
	}

	err = dec.Decode(parse)

	if err != nil {

		return err
	}

	return nil

}