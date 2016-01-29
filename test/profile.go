package main

import (
	"flag"
	"os"
	"runtime/pprof"
	"log"
	"time"
	"fmt"
)

func test() {

	time.Sleep(10000)
	fmt.Printf("%s","hello")
}

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer f.Close()
		defer pprof.StopCPUProfile()

	}

	for i:=0;i<1000000;i++ {

		test()

	}


}