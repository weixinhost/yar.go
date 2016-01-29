package main

import (
	"flag"
	"os"
	"runtime/pprof"
	"yar"
	"fmt"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func test_action(request *yar.Request,response *yar.Response) {

	fmt.Printf("%s","hello,world");
}

func main() {

	flag.Parse()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			os.Exit(-1)
		}

		pprof.StartCPUProfile(f)
		defer f.Close()
		defer pprof.StopCPUProfile()
	}

	server := yar.NewServer()

	server.RegisterHandler("test",test_action)

	server.Run()
}
