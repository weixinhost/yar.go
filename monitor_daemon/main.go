package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/weixinhost/yar.go/monitor"
)

func syncList(dataList []*monitor.MonitorData) {

	cloudWatchRequest := make(map[string]interface{})
	cloudWatchRequest["Namespace"] = "yar-server-monitor"
	metricData := make([]map[string]interface{}, 0)

	buildKey := make(map[string]*monitor.MonitorData)
	for _, v := range dataList {
		key := fmt.Sprintf("%s%d", v.Pool, v.Time)
		if _, ok := buildKey[key]; !ok {
			buildKey[key] = new(monitor.MonitorData)
			*buildKey[key] = *v
		} else {
			m := buildKey[key]
			m.SuccessTotal += v.SuccessTotal
			m.FailTotal += v.FailTotal
			m.RequestTime += v.RequestTime
		}
	}

	for _, v := range buildKey {
		t := time.Unix(int64(v.Time), 0)

		tt := v.SuccessTotal + v.FailTotal

		ht := float64(v.HostTotal)
		dt := float64(v.DownHostTotal)
		if v.HostTotal > 0 && tt > 0 {
			ht = float64(v.HostTotal) / float64(tt)
		}

		if v.DownHostTotal > 0 && tt > 0 {
			dt = float64(v.DownHostTotal) / float64(tt)
		}

		{

			item := map[string]interface{}{
				"MetricName": "ContainerTotal",
				"Unit":       "Count",
				"Value":      ht,
				"Timestamp":  t.Format(time.RFC3339Nano),
				"Dimensions": []map[string]interface{}{
					map[string]interface{}{
						"Name":  "Pool",
						"Value": v.Pool,
					},
					map[string]interface{}{
						"Name":  "State",
						"Value": "Running",
					},
				},
			}

			metricData = append(metricData, item)
		}

		{

			item := map[string]interface{}{
				"MetricName": "ContainerTotal",
				"Unit":       "Count",
				"Value":      dt,
				"Timestamp":  t.Format(time.RFC3339Nano),
				"Dimensions": []map[string]interface{}{
					map[string]interface{}{
						"Name":  "Pool",
						"Value": v.Pool,
					},
					map[string]interface{}{
						"Name":  "State",
						"Value": "Down",
					},
				},
			}

			metricData = append(metricData, item)

		}

		if v.SuccessTotal > 0 {
			item := map[string]interface{}{
				"MetricName": "RequestTotal",
				"Unit":       "Count",
				"Value":      v.SuccessTotal,
				"Timestamp":  t.Format(time.RFC3339Nano),
				"Dimensions": []map[string]interface{}{
					map[string]interface{}{
						"Name":  "Pool",
						"Value": v.Pool,
					},
					map[string]interface{}{
						"Name":  "Result",
						"Value": "Success",
					},
				},
			}

			metricData = append(metricData, item)
		}

		if v.FailTotal > 0 {
			item := map[string]interface{}{
				"MetricName": "RequestTotal",
				"Unit":       "Count",
				"Value":      v.FailTotal,
				"Timestamp":  t.Format(time.RFC3339Nano),
				"Dimensions": []map[string]interface{}{
					map[string]interface{}{
						"Name":  "Pool",
						"Value": v.Pool,
					},
					map[string]interface{}{
						"Name":  "Result",
						"Value": "Fail",
					},
				},
			}

			metricData = append(metricData, item)
		}

		if v.SuccessTotal+v.FailTotal > 0 {
			item := map[string]interface{}{
				"MetricName": "AvgRequestTime",
				"Unit":       "Milliseconds",
				"Value":      v.RequestTime / (v.SuccessTotal + v.FailTotal),
				"Timestamp":  t.Format(time.RFC3339Nano),
				"Dimensions": []map[string]interface{}{
					map[string]interface{}{
						"Name":  "Pool",
						"Value": v.Pool,
					},
				},
			}
			metricData = append(metricData, item)
		}
	}

	buildKey = make(map[string]*monitor.MonitorData)
	for _, v := range dataList {
		key := fmt.Sprintf("%s%s%d", v.Pool, v.Name, v.Time)
		if _, ok := buildKey[key]; !ok {
			buildKey[key] = new(monitor.MonitorData)
			*buildKey[key] = *v
		} else {
			m := buildKey[key]
			m.SuccessTotal += v.SuccessTotal
			m.FailTotal += v.FailTotal
			m.RequestTime += v.RequestTime
		}
	}

	for _, v := range buildKey {
		t := time.Unix(int64(v.Time), 0)

		tt := v.SuccessTotal + v.FailTotal

		ht := float64(v.HostTotal)
		dt := float64(v.DownHostTotal)
		if v.HostTotal > 0 && tt > 0 {
			ht = float64(v.HostTotal) / float64(tt)
		}

		if v.DownHostTotal > 0 && tt > 0 {
			dt = float64(v.DownHostTotal) / float64(tt)
		}

		{

			item := map[string]interface{}{
				"MetricName": "ContainerTotal",
				"Unit":       "Count",
				"Value":      ht,
				"Timestamp":  t.Format(time.RFC3339Nano),
				"Dimensions": []map[string]interface{}{
					map[string]interface{}{
						"Name":  "Pool",
						"Value": v.Pool,
					},
					map[string]interface{}{
						"Name":  "State",
						"Value": "Running",
					},
				},
			}

			metricData = append(metricData, item)
		}

		{

			item := map[string]interface{}{
				"MetricName": "ContainerTotal",
				"Unit":       "Count",
				"Value":      dt,
				"Timestamp":  t.Format(time.RFC3339Nano),
				"Dimensions": []map[string]interface{}{
					map[string]interface{}{
						"Name":  "Pool",
						"Value": v.Pool,
					},
					map[string]interface{}{
						"Name":  "State",
						"Value": "Down",
					},
				},
			}

			metricData = append(metricData, item)

		}

		if v.SuccessTotal > 0 {
			item := map[string]interface{}{
				"MetricName": "RequestTotal",
				"Unit":       "Count",
				"Value":      v.SuccessTotal,
				"Timestamp":  t.Format(time.RFC3339Nano),
				"Dimensions": []map[string]interface{}{
					map[string]interface{}{
						"Name":  "Pool",
						"Value": v.Pool,
					},
					map[string]interface{}{
						"Name":  "Name",
						"Value": v.Name,
					},
					map[string]interface{}{
						"Name":  "Result",
						"Value": "Success",
					},
				},
			}

			metricData = append(metricData, item)
		}

		if v.FailTotal > 0 {
			item := map[string]interface{}{
				"MetricName": "RequestTotal",
				"Unit":       "Count",
				"Value":      v.FailTotal,
				"Timestamp":  t.Format(time.RFC3339Nano),
				"Dimensions": []map[string]interface{}{
					map[string]interface{}{
						"Name":  "Pool",
						"Value": v.Pool,
					},
					map[string]interface{}{
						"Name":  "Name",
						"Value": v.Name,
					},
					map[string]interface{}{
						"Name":  "Result",
						"Value": "Fail",
					},
				},
			}

			metricData = append(metricData, item)
		}

		if v.SuccessTotal+v.FailTotal > 0 {
			item := map[string]interface{}{
				"MetricName": "AvgRequestTime",
				"Unit":       "Milliseconds",
				"Value":      v.RequestTime / (v.SuccessTotal + v.FailTotal),
				"Timestamp":  t.Format(time.RFC3339Nano),
				"Dimensions": []map[string]interface{}{
					map[string]interface{}{
						"Name":  "Pool",
						"Value": v.Pool,
					},
					map[string]interface{}{
						"Name":  "Name",
						"Value": v.Name,
					},
				},
			}

			metricData = append(metricData, item)
		}
	}

	cloudWatchRequest["MetricData"] = metricData
	jsonText, _ := json.Marshal(cloudWatchRequest)
	var matrics = &cloudwatch.PutMetricDataInput{}
	json.Unmarshal(jsonText, matrics)
	sess := session.New()
	sess.Config.Region = aws.String("cn-north-1")
	client := cloudwatch.New(sess)
	i := 0
	for i < len(matrics.MetricData) {
		e := i + 10
		if len(matrics.MetricData) < e {
			e = len(matrics.MetricData)
		}
		var nm = &cloudwatch.PutMetricDataInput{}
		nm.Namespace = matrics.Namespace
		nm.MetricData = matrics.MetricData[i:e]
		i += 10
		_, err := client.PutMetricData(nm)
		if err != nil {
			log.Println("PutMetricData Error:" + err.Error())
		}
	}
}

func main() {
	log.SetFlags(log.LUTC | log.LstdFlags | log.Lshortfile)
	redisHost := flag.String("redis-host", "", "Redis Host")
	flag.Parse()
	log.Println("Start Monitor Sync...")

	monitor.Setup(*redisHost, nil)
	for {
		dataList := monitor.GetLogListFromRedis(20)
		if len(dataList) < 1 {
			time.Sleep(5 * time.Second)
			continue
		}
		syncList(dataList)
	}

}
