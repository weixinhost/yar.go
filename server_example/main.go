package main

import "github.com/weixinhost/yar.go/server"

type YarClass struct{}

func (c *YarClass) Echo(str interface{}) interface{} {
	return str
}

func main() {
	//构建一个Yar Server
	yar := server.NewServer(&YarClass{})
	//基于FastHttp 构建的Server
	httpServer := server.NewHttpServer()
	//注册路由
	httpServer.RegisterHandle("/", yar)
	//启动Serve
	httpServer.Serve(":8080")
}
