#### Yar

     [Yar](https://github.com/laruence/yar) 最初是一个轻量级的PHP RPC框架.由PHP语言核心开发者[亚一程-laruence](http://www.laruence.com/)
     以PHP 扩展的方式编写.后续有了[yar-c](https://github.com/laruence/yar-c)这样纯C语言版本.
     以及基于Nginx实现的Http版本 [nginx_yar_module](https://github.com/weixinhost/nginx_yar_module)

### 特性列表

     1. 多种数据打包方案（JSON,MSGPACK）
     2. Server端非侵入式编程.与编写本地函数一致（框架使用reflect机制进行调度处理)
     3. 多种数据传输协议(TCP,UDP,HTTP)

### Example Server

```golang
package main
import ( "yar" )

func test_action(int_num int,float_num float32,str string)(string) {
	return "abcdefghj"
}

func main() {

	server := yar.NewServer("0.0.0.0", 6789)
	server.RegisterHandler("test", test_action)
	server.Run()

}

```

### Example Client

```golang

package main
import (
	"fmt"
	"yar"
)

func main(){

	client ,err := yar.NewClientWithTcp("127.0.0.1",6789)

	if err != nil {
		fmt.Printf("err:%s",err)
		return
	}

	var ret string
	err = client.Call("test",&ret,123,1.23,"abcdefg")
	fmt.Println(ret)
	fmt.Println(err)
}

```