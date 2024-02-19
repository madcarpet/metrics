package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/madcarpet/metrics/internal/collector"
)

func updateMetric(ms collector.MetricSource, i int64) {
	for {
		ms.Collect()
		time.Sleep(time.Duration(i) * time.Second)
	}
}

func sendMetric(ms collector.MetricSource, a string, i int64) {
	for {
		for k, v := range ms.GetGauge() {
			url := fmt.Sprintf("http://%s/update/gauge/%v/%v", a, k, v)
			r, err := http.Post(url, "text/plain", nil)
			if err != nil {
				panic(err)
			}
			r.Body.Close()
		}
		for k, v := range ms.GetCounter() {
			url := fmt.Sprintf("http://%s/update/counter/%v/%v", a, k, v)
			r, err := http.Post(url, "text/plain", nil)
			if err != nil {
				panic(err)
			}
			r.Body.Close()

		}
		time.Sleep(time.Duration(i) * time.Second)
	}

}

func main() {
	err := parseFlags()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	m := collector.NewMetrics()
	go updateMetric(&m, pollInterval)
	go sendMetric(&m, serverAddress, reportInterval)
	select {}
}
