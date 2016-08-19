#### Yar.go
-----
[Yar.go](https://github.com/weixinhost/yar.go) 是一个[YAR，一个轻量级的跨语言RPC框架](https://github.com/laruence/yar)的实现

#### 特性

1. 支持http客户端
2. 支持基于任意web框架部署server
3. 支持动态参数列表，但不支持默认值

-----

#### Example Server 

```go
package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/weixinhost/yar.go"
	"github.com/weixinhost/yar.go/server"
)

//这里是Yar 需要Handle的结构
type YarClass struct{}

//远程Rpc方法
func (c *YarClass) Echo() string {
	log.Println("echo handler")
	return "string"
}

func main() {
    //本示例使用了golang 自带的http 包实现http服务器。
    // 本质上，YarClient 是发送一个Post请求到相应的地址
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        //获取到完整的http request body
		body, _ := ioutil.ReadAll(r.Body)
		//得到一个server实例
        s := server.NewServer(&YarClass{})
        //注册方法,由于Golang对方法名有规定，所以提供了一个Register方法来注册方法别名。
        //当然，如果Yar客户端直接使用Echo来发起调用的话，则不需要这里注册一次
		s.Register("echo", "Echo")
		//这里接收完整的request body与一个 Writer，可以和任意web框架结合
		//数据通过Writer进行写回
		err := s.Handle(body, w)        
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

#### Example Client

```go
package main

import (
	"fmt"

	"github.com/weixinhost/yar.go"
	"github.com/weixinhost/yar.go/client"
)

func main() {
    
    //初始化一个客户端。
    //目前仅支持http,https 的Yar服务端
	client, err := client.NewClient("http://127.0.0.1:8080")
    
	if err != nil {
		fmt.Println("error", err)
	}

	//这是默认值
	client.Opt.Timeout = 1000 * 30 //30s
	//这是默认值，目前只支持json
	client.Opt.Packager = "json"
    client.Opt.DNSCache = true \\开启DNS缓存,一旦开启DNS缓存，将会使用内存进行60秒的内存缓存设置。
    
    //定义Yar的服务端方法返回值
	var ret interface{}

	callErr := client.Call("echo", &ret)
    
    //错误判断
	if callErr != nil {
		fmt.Println("error", callErr)
	}
    
	fmt.Println("data", ret)
}

```



