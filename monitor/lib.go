package monitor

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	redis "gopkg.in/redis.v3"
)

type MonitorData struct {
	Pool          string
	Name          string
	Provider      string
	RequestTime   int
	IsSuccess     bool
	Time          int
	SuccessTotal  int
	FailTotal     int
	HostTotal     int
	DownHostTotal int
}

type RealTimeMonitorHandle func(pool, name, addr string, msg string)

var redisLstKey string = "__Yar_Monitor_CacheList__"
var redisClient *redis.Client
var cacheMutex sync.Mutex
var cacheList map[string]*MonitorData
var realTimeMonitor RealTimeMonitorHandle

func init() {
	cacheList = make(map[string]*MonitorData)
}

func Setup(redisHost string, h RealTimeMonitorHandle) {
	realTimeMonitor = h
	if redisClient != nil {
		redisClient.Close()
	}

	opt := &redis.Options{}
	opt.Addr = redisHost
	opt.DB = 7
	opt.IdleTimeout = 60 * time.Second
	opt.WriteTimeout = 10 * time.Second
	opt.ReadTimeout = 10 * time.Second
	opt.MaxRetries = 3
	redisClient = redis.NewClient(opt)

	go SyncLogToRedis()
}

func SetServiceMonitor(pool, name, provider string, requestTime int, healthHostTotal int, downHostTotal int, isSuccess bool) {

	if redisClient == nil {
		return
	}

	t := int(time.Now().Unix())
	t = t - (t % 30)

	monitor := &MonitorData{
		Pool:          pool,
		Name:          name,
		Provider:      provider,
		RequestTime:   requestTime,
		IsSuccess:     isSuccess,
		HostTotal:     healthHostTotal,
		DownHostTotal: downHostTotal,
		Time:          t,
	}

	if isSuccess {
		monitor.SuccessTotal = 1
	} else {
		monitor.FailTotal = 1
	}

	key := fmt.Sprintf("%s%s%s%d", pool, name, provider, t)

	cacheMutex.Lock()
	if _, ok := cacheList[key]; !ok {
		cacheList[key] = monitor
	} else {
		m := cacheList[key]
		m.RequestTime += requestTime
		if isSuccess {
			m.SuccessTotal++
		} else {
			m.FailTotal++
		}
	}
	cacheMutex.Unlock()
}

func SyncLogToRedis() {

	for {
		if len(cacheList) < 1 {
			time.Sleep(5 * time.Second)
			continue
		}
		cacheMutex.Lock()
		list := cacheList
		cacheList = make(map[string]*MonitorData)
		cacheMutex.Unlock()
		for _, v := range list {
			data, err := json.Marshal(v)
			if err != nil {
				continue
			}
			redisClient.LPush(redisLstKey, string(data))
		}
		time.Sleep(2 * time.Second)
	}
}

func GetLogListFromRedis(max int) []*MonitorData {

	var lst []*MonitorData

	for i := 0; i < max; i++ {
		cmd := redisClient.LPop(redisLstKey)

		if cmd.Err() != nil {
			break
		}

		data := cmd.Val()

		m := new(MonitorData)

		err := json.Unmarshal([]byte(data), m)

		if err == nil {
			lst = append(lst, m)
		}
	}
	return lst
}

func RealTimeMonitor(pool, name, addr string, msg string) {

	if realTimeMonitor == nil || redisClient == nil {
		return
	}

	key := fmt.Sprintf("__Yar_Monitor_RealTime__:%s:%s:%s", pool, name, addr)
	cmd := redisClient.Get(key)
	t, _ := cmd.Int64()
	now := time.Now().Unix()
	if now-t < 60 {
		return
	}
	redisClient.Set(key, now, 15*time.Minute)
	go realTimeMonitor(pool, name, addr, msg)
}
