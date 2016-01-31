package yar
import (
	"yar/transports"
)
// 定义RPC请求的函数,可以为任意类型的函数.但是具备以下限制:\n
// 1. 无论是接收的值或是返回的值都需要能够正确的被打包
// 2. 不支持多返回值
type Handler interface{}

// 定义RPC请求的过滤器原型.
// 由于在真实的用于rpc调用的函数是非侵入式的,目前无法直接在函数内获取到一些rpc的信息.因此,可以注册一系列的过滤器来完成rpc相关
// 信息的探测.当需要拦截当前请求的时候,可以返回false.后续操作将被直接拦截
type ServerFilter func(server *Server,conn transports.TransportConnection,request *Request,response *Response)(bool)

