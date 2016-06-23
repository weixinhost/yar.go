package yar

import "fmt"

type ErrorEnum int

const (
	//网络错误
	ErrorNetwork ErrorEnum = 1
	//配置错误
	ErrorConfig ErrorEnum = 2
	//请求参数错误
	ErrorParam ErrorEnum = 3
	//数据包打包/解包错误
	ErrorPackager ErrorEnum = 4
	//协议解析错误
	ErrorProtocol ErrorEnum = 5
	//请求验证错误
	ErrorVerify ErrorEnum = 6
	//加密解密错误
	ErrorEncrypt ErrorEnum = 7
	//返回数据错误
	ErrorResponse ErrorEnum = 8
	//返回数据错误
	ErrorRequest ErrorEnum = 9
)

func (e ErrorEnum) String() string {

	switch e {
	case ErrorNetwork:
		return "NetWork Error"
	case ErrorConfig:
		return "Config Error"
	case ErrorParam:
		return "Param Error"
	case ErrorPackager:
		return "Packager Error"
	case ErrorProtocol:
		return "Protocol Error"
	case ErrorVerify:
		return "Verify Error"
	case ErrorEncrypt:
		return "Encrypt Error"
	case ErrorResponse:
		return "Response Error"
	case ErrorRequest:
		return "Request Error"
	}

	return "Unknow Error"
}

type Error struct {
	t ErrorEnum
	m string
}

func NewError(t ErrorEnum, m string) *Error {
	return &Error{t: t, m: m}
}

func (e *Error) String() string {
	return fmt.Sprintf("[%s] %s", e.t, e.m)
}

func (e *Error) Assert(t ErrorEnum) bool {
	return e.t == t
}
