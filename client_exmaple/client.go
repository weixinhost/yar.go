package main

import (
	"fmt"

	yar "github.com/weixinhost/yar.go"
	"github.com/weixinhost/yar.go/client"
)

func main() {

	client, err := client.NewClient("http://127.0.0.1:8080")

	if err != nil {
		fmt.Println("error", err)
	}

	//这是默认值
	client.Opt.Timeout = 1000 * 30 //30s
	//这是默认值
	client.Opt.Packager = "json"
	//这是默认值
	client.Opt.Encrypt = false
	//这是默认值
	client.Opt.EncryptPrivateKey = ""
	//这是默认值
	client.Opt.MagicNumber = yar.MagicNumber

	//	param := 1

	var ret interface{}

	callErr := client.Call("echo", &ret)

	if callErr != nil {
		fmt.Println("error", callErr)
	}

	fmt.Println("data", ret)
}
