package packager

import (
	"encoding/json"
	"strings"
)

func JsonPack(v interface{}) ([]byte, error) {
	data, err := json.Marshal(v)
	return data, err
}

func JsonUnpack(data []byte, v interface{}) error {
	d := json.NewDecoder(strings.NewReader(string(data)))
	d.UseNumber()
	err := d.Decode(v)
	return err
}
