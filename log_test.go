package yar

import (
	"testing"
)

func TestLog(t *testing.T) {

	log,err := LogNew("~/a.log")

	if(err != nil){

		t.Log(err.Error())

		return
	}

	log.Normal("%s %d","string",1)

	t.Log(log.buffers)

}