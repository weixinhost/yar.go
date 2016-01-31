package main
import (

	"runtime"
	//"yar"
	"net/http"
	"fmt"
)


type Custom struct {

}

func(c *Custom)ServeHTTP(writer http.ResponseWriter,request *http.Request){

	fmt.Printf("sss")
}

func test_action(int_num int,float_num float32,str string)(string) {

	return "abcdefghj"

}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	/*
	server,_:= yar.NewServer("udp",":6790")
	server.RegisterHandler("test", test_action)
	server.Serve()
	*/

	http.HandleFunc("/",func(writer http.ResponseWriter,request *http.Request){

		fmt.Printf("%s %s","start",request.Body)
	})

	http.ListenAndServe(":8080",nil)

}