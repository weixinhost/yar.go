package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/weixinhost/yar.go"
	"github.com/weixinhost/yar.go/server"
)

type YarClass struct{}

func (c *YarClass) Echo() string {
	log.Println("echo handler")
	return "string"
}

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		s := server.NewServer(&YarClass{})

		s.Opt.MagicNumber = yar.MagicNumber
		s.Register("echo", "Echo")
		//这里接收完整的request body与一个 Writer，可以和任意web框架结合
		//数据通过Writer进行写回
		s.Handle(body, w)
	})
	log.Fatal(http.ListenAndServe(":8080", nil))

}
