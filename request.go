package yar

type Request struct {
	Protocol *Protocol   `json:"-" msgpack:"-"`
	Id       uint32      `json:"i" msgpack:"i"`
	Method   string      `json:"m" msgpack:"m"`
	Params   interface{} `json:"p" msgpack:"p"`
}


func NewRequest() (request *Request){

	request = new(Request)
	return request
}

