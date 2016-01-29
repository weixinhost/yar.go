package main
import (
	"yar"
	"fmt"
)

func test_action(request *yar.Request, response *yar.Response) {

	fmt.Printf("%s", "hello,world")
	response.Return("abcdefgh")
}

func main() {

	server := yar.NewServer("0.0.0.0", 6789)
	server.RegisterHandler("test", test_action)
	server.Run()

}
