package client

import "errors"
import "strings"

var supportNets = []string{
	"http",
	"https",
	"tcp",
	"udp",
	"unix",
}

func parseAddrNetName(addr string) (string, error) {

	splitIndex := strings.IndexAny(addr, ":")

	if splitIndex < 3 {
		return "", errors.New("parse addr error")
	}

	protocol := addr[0:splitIndex]

	for _, item := range supportNets {
		if strings.ToLower(protocol) == item {
			return item, nil
		}
	}

	return "", errors.New("unsupport net protocol: " + protocol)

}
