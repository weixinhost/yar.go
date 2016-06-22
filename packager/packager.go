package packager

import (
	"errors"
	"strings"
	//"bytes"
	"bytes"
)

type PackFunc func(v interface{}) ([]byte, error)

type UnpackFunc func(data []byte, v interface{}) error

func Pack(name []byte, v interface{}) ([]byte, error) {

	s := strings.ToLower(bytes.NewBuffer(name).String())

	if strings.Contains(s, "json") {

		return JsonPack(v)
	}

	return nil, errors.New("unsupported packager")

}

func Unpack(name []byte, data []byte, v interface{}) error {

	s := strings.ToLower(bytes.NewBuffer(name).String())

	if strings.Contains(s, "json") {

		return JsonUnpack(data, v)

	}

	return errors.New("unsupported packager")
}
