package yar

type Response struct {
	Protocol *Header     `json:"-" msgpack:"-"`
	Id       uint32      `json:"i" msgpack:"i"`
	Error    string      `json:"e" msgpack:"e"`
	Out      string      `json:"o" msgpack:"o"`
	Status   ErrorType   `json:"s" msgpack:"s"`
	Retval   interface{} `json:"r" msgpack:"r"`
}

func NewResponse() (response *Response) {

	response = new(Response)

	return response
}

func (self *Response) Exception(msg string) {

	self.Status = ERR_OUTPUT
	self.Error = msg
}

func (self *Response) Output(msg string) {
	self.Out += msg
}

func (self *Response) Return(v interface{}) (err error) {
	self.Retval = v
	return nil
}
