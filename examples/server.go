package main
import (

	"runtime"
	"net/rpc"
)

func test_action(int_num int,float_num float32,str string)(string) {

	return "abcdefghj"

}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	server :=rpc.NewServer()

	server.Register(test_action)

	server.Accept()
}