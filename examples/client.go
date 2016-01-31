package main
import (
	"fmt"
	"yar"

)

type UUidResponse struct {

	Uuid float64		`json:"uuid"`
}


func main(){

	client := yar.NewClient("http","http://127.0.0.1:6790")

	ret := new(UUidResponse)

	err := client.Call("uuid",&ret)

	if err != nil {

		fmt.Println(err)

		return ;
	}

	fmt.Println(ret)

}

