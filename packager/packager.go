package packager

import (
	"errors"
	"strings"
	//"bytes"
	"bytes"
)

type PackFunc func(v interface{}) ([]byte, error)

type UnpackFunc func(data []byte, v interface{}) error

func Pack(name []byte ,v interface{}) ([]byte, error) {

	switch strings.ToLower(bytes.NewBuffer(name).String()) {

	case "json":
		{
			return JsonPack(v)
		}
		break
	case "msgpack":
		{
		}
		break
	}

	return nil, errors.New("unsupported packager")
}

func Unpack(name []byte, data []byte, v interface{}) error {

	switch strings.ToLower(bytes.NewBuffer(name).String()) {

	case "json":
		{
			return JsonUnpack(data, v)
		}
		break

	case "msgpack":
		{
		}
		break
	}

	return errors.New("unsupported packager")
}
