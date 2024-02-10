package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/madcarpet/metrics/internal/collector"
)

type metricSource interface {
	Collect()
	GetGauge() map[string]float64
	GetCounter() map[string]int64
}

func updateMetric(ms metricSource) {
	for {
		ms.Collect()
		time.Sleep(2 * time.Second)
	}
}

func sendMetric(ms metricSource) {
	for {
		for k, v := range ms.GetGauge() {
			url := fmt.Sprintf("http://localhost:8080/update/gauge/%v/%v", k, v)
			r, _ := http.Post(url, "text/plain", nil)
			fmt.Println(k, v)
			fmt.Println("This is gauge response: ", r.StatusCode)
		}
		for k, v := range ms.GetCounter() {
			url := fmt.Sprintf("http://localhost:8080/update/counter/%v/%v", k, v)
			r, _ := http.Post(url, "text/plain", nil)
			fmt.Println(k, v)
			fmt.Println("This is counter response: ", r.StatusCode)
		}
		time.Sleep(10 * time.Second)
	}

}

func main() {
	m := collector.NewMetrics()
	go updateMetric(&m)
	go sendMetric(&m)
	select {}
}
