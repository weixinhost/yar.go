package yar
import (
	"math/rand"
	"time"
)

type Request struct {
	Protocol *Protocol   	`json:"-" msgpack:"-"`
	Id       uint32      	`json:"i" msgpack:"i"`
	Method   string      	`json:"m" msgpack:"m"`
	Params   interface{} 	`json:"p" msgpack:"p"`
}

func NewRequest() (request *Request){
	request = new(Request)
	rand.Seed(time.Now().Unix())
	request.Id = rand.Uint32()
	return request
}

