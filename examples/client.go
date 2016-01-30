package main
import (
	"fmt"
	"yar"
)

func main(){

	client ,err := yar.NewClientWithTcp("127.0.0.1",6790)

	if err != nil {
		fmt.Printf("err:%s",err)
		return
	}

	var ret map[string]uint64
	err = client.Call("uuid",&ret,1)
	fmt.Println(ret)
	fmt.Println(err)
}

