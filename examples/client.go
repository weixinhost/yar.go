package main
import (
	"fmt"
	"yar"

)

type UUidResponse struct {

	Uuid float64		`json:"uuid"`
}


func main(){

	client := yar.NewClient("unix","/tmp/a3.sock")

	ret := new(UUidResponse)

	err := client.Call("uuid",&ret)

	if err != nil {

		fmt.Println(err)

		return ;
	}

	fmt.Println(ret)

}

