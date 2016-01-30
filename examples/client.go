package main
import (
	"fmt"
	"yar"
)

type UUidResponse struct {

	Uuid float64		`json:"uuid"`
}

func main(){

	client ,err := yar.NewClientWithTcp("127.0.0.1",6790)

	if err != nil {
		fmt.Printf("err:%s",err)
		return
	}

	ret := new(UUidResponse)

	err = client.Call("uuid",&ret)

	if err != nil {

		fmt.Println(err)

		return ;
	}

	fmt.Println(ret)

}

