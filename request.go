package yar

type Request struct {

	protocol 	*Protocol

	id 		    uint32

	method 		string

	params 		interface{}

}

func RequestNew() *Request {

	request := new(Request)

	return request;

}



