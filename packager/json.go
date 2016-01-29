package packager

import "encoding/json"

func JsonPack(v interface{}) ([]byte,error){

	data,err := json.Marshal(v)

	return data,err
}

func JsonUnpack(data []byte,v interface{}) error{

	err := json.Unmarshal(data,v)

	return err
}
