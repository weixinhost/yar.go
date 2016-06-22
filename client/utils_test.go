package client

import "testing"

func TestParseAddrNet(t *testing.T) {

	tests := map[string]string{
		"http":  "http://test.local.com",
		"https": "https://test.local.com",
		"tcp":   "tcp://127.0.0.1",
		"udp":   "udp://abcd",
	}

	for k, v := range tests {

		n, e := parseAddrNetName(v)

		if e != nil {
			t.Fatal(e)
		}
		if n != k {
			t.Fatal(v, n)
		}

	}

}
