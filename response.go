package yar

type Response struct {
	Protocol *Protocol  	`json:"-" msgpack:"-"`
	Id     int64       		`json:"i" msgpack:"i"`
	Error  string      		`json:"e" msgpack:"e"`
	Output string      		`json:"o" msgpack:"o"`
	Status ExceptionType    `json:"s" msgpack:"s"`
	Retval []byte 			`json:"r" msgpack:"r"`
}
