package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"time"
)

func test() {

	time.Sleep(10000)
	fmt.Printf("%s", "hello")
}

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {


}
