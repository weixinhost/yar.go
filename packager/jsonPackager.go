package yar

import (
	"encoding/json"
	"bytes"
)

type JsonPackager struct {
	
	Packager
}

func JsonPackagerNew() *JsonPackager {
	
	packager := new(JsonPackager)
	
	return packager
}


func (self *JsonPackager)Pack(data map[string]interface{}) *bytes.Buffer {

	var ret,err = json.Marshal(data)

	if(err != nil) {

		return nil

	}

	return bytes.NewBuffer(ret)

}

func (self *JsonPackager)Unpack(data *bytes.Buffer) map[string]interface{} {
	
	var ret map[string]interface{}
	
	json.Unmarshal(data.Bytes(),&ret)
	
	return ret
	
}




