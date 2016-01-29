package main
import (
	"yar"
	"runtime/pprof"
	"flag"
	"os"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main(){

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
	server.Run()
}
