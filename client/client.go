package client

import (
	"bytes"
	"crypto/tls"
	"errors"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"

	yar "github.com/weixinhost/yar.go"
	"github.com/weixinhost/yar.go/packager"
	"github.com/weixinhost/yar.go/transports"
)

type Client struct {
	hostname  string
	net       string
	transport transports.Transport
	Opt       *yar.Opt
}

// 获取一个YAR 客户端
// addr为带请求协议的地址。支持以下格式
// http://xxxxxxxx
// https://xxxx.xx.xx
// tcp://xxxx
// udp://xxxx
func NewClient(addr string) (*Client, *yar.Error) {
	netName, err := parseAddrNetName(addr)
	if err != nil {
		return nil, yar.NewError(yar.ErrorParam, err.Error())
	}

	client := new(Client)

	client.hostname = addr
	client.net = netName
	client.Opt = yar.NewOpt()
	client.init()
	return client, nil
}

func (client *Client) init() {
	switch client.net {
	case "tcp", "udp", "unix":
		{
			client.transport, _ = transports.NewSock(client.net, client.hostname)
			break
		}
	}

}

func (client *Client) Call(method string, ret interface{}, params ...interface{}) *yar.Error {

	if client.net == "http" || client.net == "https" {
		return client.httpHandler(method, ret, params...)
	}

	return yar.NewError(yar.ErrorConfig, "unsupported non http protocol")

}

func (client *Client) initRequest(method string, params ...interface{}) (*yar.Request, *yar.Error) {

	r := yar.NewRequest()

	if len(method) < 1 {
		return nil, yar.NewError(yar.ErrorParam, "call empty method")
	}

	if params == nil {
		r.Params = []interface{}{}
	} else {
		r.Params = params
	}

	r.Method = method

	r.Protocol.MagicNumber = client.Opt.MagicNumber
	r.Protocol.Id = r.Id
	return r, nil
}

func (client *Client) packRequest(r *yar.Request) ([]byte, *yar.Error) {

	var sendPackager []byte
	packagerName := client.Opt.Packager

	if len(packagerName) < yar.PackagerLength {
		sendPackager = []byte(packagerName)
	} else {
		sendPackager = []byte(packagerName[0:yar.PackagerLength])
	}

	var p [8]byte
	for i, s := range sendPackager {
		p[i] = s
	}

	r.Protocol.Packager = p
	pack, err := packager.Pack(sendPackager, r)

	if err != nil {
		return nil, yar.NewError(yar.ErrorPackager, err.Error())
	}

	return pack, nil
}

func (client *Client) readResponse(reader io.Reader, ret interface{}) *yar.Error {

	allBody, err := ioutil.ReadAll(reader)

	if err != nil {
		return yar.NewError(yar.ErrorResponse, "Read Response Error:"+err.Error())
	}

	if len(allBody) < (yar.ProtocolLength + yar.PackagerLength) {
		return yar.NewError(yar.ErrorResponse, "Response Parse Error:"+string(allBody))
	}

	protocolBuffer := allBody[0 : yar.ProtocolLength+yar.PackagerLength]

	protocol := yar.NewHeader()

	protocol.Init(bytes.NewBuffer(protocolBuffer))

	bodyLength := protocol.BodyLength - yar.PackagerLength

	if uint32(len(allBody)-(yar.ProtocolLength+yar.PackagerLength)) < uint32(bodyLength) {
		return yar.NewError(yar.ErrorResponse, "Response Content Error:"+string(allBody))
	}

	bodyBuffer := allBody[yar.ProtocolLength+yar.PackagerLength:]

	response := new(yar.Response)
	err = packager.Unpack([]byte(client.Opt.Packager), bodyBuffer, &response)

	if err != nil {
		return yar.NewError(yar.ErrorPackager, "Unpack Error:"+err.Error())
	}

	if response.Status != yar.ERR_OKEY {
		return yar.NewError(yar.ErrorResponse, response.Error)
	}

	if ret != nil {

		packData, err := packager.Pack([]byte(client.Opt.Packager), response.Retval)

		if err != nil {
			return yar.NewError(yar.ErrorPackager, "pack response retval error:"+err.Error())
		}

		err = packager.Unpack([]byte(client.Opt.Packager), packData, ret)

		if err != nil {
			return yar.NewError(yar.ErrorPackager, "pack response retval error:"+err.Error())
		}
	}

	return nil
}

func (client *Client) httpHandler(method string, ret interface{}, params ...interface{}) *yar.Error {

	r, err := client.initRequest(method, params...)

	if err != nil {
		return err
	}

	packBody, err := client.packRequest(r)

	if err != nil {
		return err
	}

	r.Protocol.BodyLength = uint32(len(packBody) + yar.PackagerLength)

	postBuffer := bytes.NewBuffer(r.Protocol.Bytes().Bytes())
	postBuffer.Write(packBody)

	//todo 停止验证HTTPS请求
	tr := &http.Transport{
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		MaxIdleConns:        32,
		MaxIdleConnsPerHost: 32,
	}

	if client.Opt.DNSCache == true {
		tr.Dial = func(network string, address string) (net.Conn, error) {
			separator := strings.LastIndex(address, ":")
			ips, err := globalResolver.Lookup(address[:separator])
			if err != nil {
				return nil, errors.New("Lookup Error:" + err.Error())
			}
			if len(ips) < 1 {
				return nil, errors.New("Lookup Error: No IP Resolver Result Found")
			}
			return net.Dial("tcp", ips[0].String()+address[separator:])
		}
	}

	httpClient := &http.Client{
		Transport: tr,
		Timeout:   time.Duration(client.Opt.Timeout) * time.Millisecond,
	}

	resp, postErr := httpClient.Post(client.hostname, "application/json", postBuffer)

	if postErr != nil {
		return yar.NewError(yar.ErrorNetwork, postErr.Error())
	}

	responseErr := client.readResponse(resp.Body, ret)
	return responseErr
}

func (client *Client) sockHandler(method string, ret interface{}, params ...interface{}) *yar.Error {
	return yar.NewError(yar.ErrorParam, "unsupported sock request")
}
