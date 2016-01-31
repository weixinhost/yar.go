#### Yar

[Yar](https://github.com/laruence/yar) 最初是一个轻量级的PHP RPC框架.由PHP语言核心开发者[亚一程-laruence](http://www.laruence.com/)
以PHP 扩展的方式编写.后续有了[yar-c](https://github.com/laruence/yar-c)这样纯C语言版本.
以及基于Nginx实现的Http版本 [nginx_yar_module](https://github.com/weixinhost/nginx_yar_module)

### 特性列表

     1. 多种数据打包方案（JSON,MSGPACK）
     2. Server端非侵入式编程.与编写本地函数一致（框架使用reflect机制进行调度处理)
     3. 多种数据传输协议(TCP,UDP,HTTP)

### Example Server

```go

package main
import ( "yar" )


func test_action(int_num int,float_num float32,str string)(string) {
	return "abcdefghj"
}

func main(){
    server,_:= yar.NewServer("tcp",":6790")
	server.RegisterHandler("echo", test_action)
	//进入服务器事件循环
    server.Serve()
}

```

### Example Client

```go

package main
import (
	"fmt"
	"yar"
)

func main(){

    client := yar.NewClient("tcp","127.0.0.1:6790")
    var ret string
    err := client.Call("echo",&ret,"hello")
}

```

### Documention

    可以使用godoc查看开发文档
