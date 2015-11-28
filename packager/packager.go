package yar

import (
	"strings"
	"bytes"
	"errors"
)

type Packager interface {

	Pack(data map[string]interface{}) *bytes.Buffer
	Unpack(data *bytes.Buffer)  map[string]interface{}

}


func PackagerNew(packagerName string)  ( Packager, error){
		
	switch strings.ToLower(packagerName)  {
	
		case "json" : 
			
			ret := JsonPackagerNew()
			
			var err error = nil

			return  ret,err
		
		break
		
		case "msgpack" :
		break	
	}
	
	var ret Packager = nil
	
	var err error = errors.New("Unsupported Packager")

	return ret,err
	 
}





