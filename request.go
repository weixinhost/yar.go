package yar
import "fmt"



type Request struct {
	Protocol *Protocol   `json:"-" msgpack:"-"`
	Id       uint32      `json:"i" msgpack:"i"`
	Method   string      `json:"m" msgpack:"m"`
	Params   interface{} `json:"p" msgpack:"p"`
}




func (self *Request)GetParamWithInt32(index int)(ret int32){

	fmt.Println(self.Params)

	return int32(self.Params.([]interface{})[index].(float64))

}

func (self *Request)GetParamWithInt64(index int)(ret int64){

	return int64(self.Params.([]interface{})[index].(float64))

}

func (self *Request)GetParamWithBool(index int)(ret bool){

	return bool(self.Params.([]interface{})[index].(bool))

}


func (self *Request)GetParamWithString(index int)(ret string){

	return string(self.Params.([]interface{})[index].(string))

}

func (self *Request)GetParamWithArray(index int)(ret []interface{}) {

	return self.Params.([]interface{})[index].([]interface{})

}


func (self *Request)GetParamWithMap(index int)(ret map[string]interface{}) {

	return self.Params.([]interface{})[index].(map[string]interface{})

}

