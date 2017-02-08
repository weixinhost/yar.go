package main

import (
	"flag"
	"log"
	"time"

	"math/rand"

	"github.com/weixinhost/yar.go/host_sync"
)

func main() {

	log.SetFlags(log.LUTC | log.LstdFlags | log.Lshortfile)

	dockerAPI := flag.String("docker-api", "", "Docker API")
	redisHost := flag.String("redis-host", "", "Redis Host")
	flag.Parse()

	host_sync.SetRedisHost(*redisHost)
	host_sync.SetDockerAPI(*dockerAPI)

	log.Println("Start Host Sync...")

	for {
		rand.Seed(time.Now().Unix())
		err := host_sync.SyncAllHostList()
		if err != nil {
			log.Println(err)
		}
		time.Sleep(time.Duration((5 + rand.Intn(3))) * time.Second)
	}
}
