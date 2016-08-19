package client

import (
	"fmt"
	"testing"
	"time"
)

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

func TestDNS(t *testing.T) {

	resolver := NewResolver(10, 10*time.Second)

	var list []string = []string{
		"www.sina.com.cn",
		"www.baidu.com",
		"127.0.0.1:5600",
		"www.google.com",
		"www.facebook.com",
		"www.weixinhost.com",
		"follower.services-azure.weixinhost.com",
		"follower.service.weixinhost.com",
		"test.notfound.weixinhost.com",
		"alpwosk91712as12212kao1",
		"localhost",
		"asf234232e-10238121.com",
		"https://www.baidu.com",
	}

	for _, v := range list {
		r, err := resolver.Lookup(v)
		fmt.Println(v, r, err)
	}
	fmt.Println("")
	for _, v := range list {
		r, err := resolver.Lookup(v)
		fmt.Println(v, r, err)
	}

	time.Sleep(10 * time.Second)
	fmt.Println("")
	for _, v := range list {
		r, err := resolver.Lookup(v)
		fmt.Println(v, r, err)
	}
}
