package yar

import (
	"bytes"
	"errors"
	"math/rand"
	"yar/packager"
	"yar/transports"
	"encoding/gob"
	"strings"
	"net/http"
)

// Client的配置项
type Opt int

const (
	CONNECTION_TIMEOUT Opt = 1	//连接超时
	TIMEOUT            Opt = 2	//整体超时
	PACKAGER           Opt = 3	//打包协议.目前支持 "json"
)

const (
	DEFAULT_PACKAGER           			= "json"		// 默认打包协议
	DEFAULT_TIMEOUT_SECOND            = 5000			// 默认超时.包含连接超时.因此,rpc函数的执行超时为 TIMEOUT - CONNECTION_TIMEOUT
	DEFAULT_CONNECTION_TIMEOUT_SECOND = 1000			// 默认链接超时
)

//用于yar请求的客户端
type Client struct {
	net string			//网络传输协议.支持 "tcp","udp","http","unix"等值
	hostname string		//用于初始化网络链接的信息,入 ip:port domain:port 等
	request   *Request	//请求体
	transport transports.Transport
	opt       map[Opt]interface{}	//配置项
}

//初始化一个客户端
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

	case "tcp","udp","unix" : {
		client.transport,_ = transports.NewSock(client.net,client.hostname)
		break
	}
	}

}

func (self *Client) initOpt() {

	self.opt[CONNECTION_TIMEOUT] = DEFAULT_CONNECTION_TIMEOUT_SECOND
	self.opt[TIMEOUT] = DEFAULT_TIMEOUT_SECOND
	self.opt[PACKAGER] = DEFAULT_PACKAGER

}

//配置项操作
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

func (self *Client) sockCall(method string,ret interface{},params ...interface{}) (err error) {

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
	//这里需要优化,需要干掉这次pack/unpack
	pack_data,err := packager.Pack(self.request.Protocol.Packager[:],response.Retval)
	err = packager.Unpack(self.request.Protocol.Packager[:],pack_data,ret)

	return err
}

func (self *Client) httpCall(method string,ret interface{},params ...interface{}) (err error) {

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

	post_buffer := bytes.NewBuffer(self.request.Protocol.Bytes().Bytes())
	post_buffer.Write(pack)
	resp,err := http.Post(self.hostname,"application/json",post_buffer)
	if err != nil {

		return err
	}

	protocol_buffer := make([]byte, PROTOCOL_LENGTH + PACKAGER_LENGTH)
	resp.Body.Read(protocol_buffer)
	self.request.Protocol.Init(bytes.NewBuffer(protocol_buffer))
	body_buffer := make([]byte, self.request.Protocol.BodyLength - PACKAGER_LENGTH)
	resp.Body.Read(body_buffer)

	response := new(Response)
	err = packager.Unpack([]byte(self.opt[PACKAGER].(string)), body_buffer, &response)

	if response.Status != ERR_OKEY {
		return errors.New(response.Error)
	}

	//这里需要优化,需要干掉这次pack/unpack
	pack_data,err := packager.Pack(self.request.Protocol.Packager[:],response.Retval)
	err = packager.Unpack(self.request.Protocol.Packager[:],pack_data,ret)
	return err
}

//执行一次rpc请求.
//method为请求的方法名.ret参数必须是一个指针类型,用于接收rpc结果.params为rpc函数的形参列表
func (self *Client) Call(method string, ret interface{},params ...interface{}) (err error) {

	switch self.net {

	case "tcp" , "udp" , "unix":
		{
			return self.sockCall(method, ret,params...)
		}

	case "http" :{
		return self.httpCall(method,ret,params...)
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