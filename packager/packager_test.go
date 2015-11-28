package yar

import (

	"bytes"
	"testing"
)

func TestPackJson(t *testing.T){

	packager,err :=  PackagerNew("json")

	if(err != nil){

		t.Errorf("%s",err.Error())

		return

	}

	var data = make(map[string]interface{})

	data["a"] = "abc";

	data["b"] = 1;

	ret := packager.Pack(data)

	t.Log(ret)





}


func TestUnPackJson(t *testing.T) {

	packager,err :=  PackagerNew("json")

	if(err != nil){

		t.Errorf("%s",err.Error())

		return

	}

	json_string := "{\"a\":{\"b\":[1,2,3]}}"

	ret := packager.Unpack(bytes.NewBufferString(json_string))

	t.Log(ret)


}


func main() {
	

}