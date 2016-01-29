package yar
import (
	"yar/packager"
	"bytes"
)

type Response struct {
	Protocol *Protocol     	`json:"-" msgpack:"-"`
	Id       uint32         	`json:"i" msgpack:"i"`
	Error    string        	`json:"e" msgpack:"e"`
	Out   	 string        	`json:"o" msgpack:"o"`
	Status   ErrorType 		`json:"s" msgpack:"s"`
	Retval   string        	`json:"r" msgpack:"r"`
}

func (self *Response) Exception(msg string) {

	self.Status = ERR_OUTPUT
	self.Error 	= msg

}


func (self *Response) Output(msg string){

	self.Out += msg

}


func (self *Response) Return(v interface{}) (err error) {

	pack,err := packager.Pack(self.Protocol.Packager[0:],v)

	if err != nil {

		return err
	}
	self.Retval = bytes.NewBuffer(pack).String()

	return nil

}
