package yar_test
import (
	//"testing"
	"yar"
	"runtime"
)

func test_action() string {

	return "world!"
}

func ExampleServer(){

	//设置服务器最大使用的CPU数量
	runtime.GOMAXPROCS(runtime.NumCPU())

	server,_:= yar.NewServer("tcp",":6790")
	server.RegisterHandler("echo", test_action)

	//进入服务器事件循环
	server.Serve()

}

func ExampleClient() {

	client := yar.NewClient("tcp","127.0.0.1:6790")
	var ret string
	err := client.Call("echo",&ret,"hello")

	if err != nil {

	}

	//Output: world!
}
