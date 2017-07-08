package main

import "github.com/weixinhost/yar.go/server"

type YarClass struct{}

func (c *YarClass) Echo(str interface{}) interface{} {
	return str
}

func main() {
	//构建一个Yar Server

	//基于FastHttp 构建的Server
	httpServer := server.NewHttpServer()
	//注册路由
	httpServer.RegisterHandle("/", func() *server.Server {
		server := server.NewServer(&YarClass{})
		server.Register("echo", "Echo")
		return server
	})
	//启动Serve
	httpServer.Serve(":8080")
}
