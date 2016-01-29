package packager
import (
	"strings"
	"errors"
)

type PackFunc func (v interface{}) ([]byte,error)

type UnpackFunc func (data []byte,v interface{}) error


func Pack(name string,v interface {}) ([]byte,error) {

	switch strings.ToLower(name) {

	case "json" : {

		return JsonPack(v)

	}
		break;
	case "msgpack" : {}
		break;
	}

	return nil,errors.New("unsupported packager")
}


func Unpack(name string,data []byte,v interface{}) error {

	switch strings.ToLower(name) {

	case "json" : {

		return JsonUnpack(data,v)

	}
		break;

	case "msgpack":{}
		break;
	}

	return errors.New("unsupported packager")
}

