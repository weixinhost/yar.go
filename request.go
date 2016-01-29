package yar

type Request struct {
	Protocol *Protocol   `json:"-" msgpack:"-"`
	Id       int64       `json:"i" msgpack:"i"`
	Method   string      `json:"m" msgpack:"m"`
	Params   interface{} `json:"p" msgpack:"p"`
}
