package main

import (
	"fmt"
	"os"
	"time"

	"github.com/madcarpet/metrics/internal/adapter/http"
	"github.com/madcarpet/metrics/internal/adapter/storage/memstorage"
	"github.com/madcarpet/metrics/internal/service/metrics"
)

type reporter interface {
	ReportMetrics() error
}

type collectService interface {
	Collect(ms []string) error
}

func metricCollecting(pi int64, c collectService, ms []string) {
	for {
		c.Collect(ms)
		time.Sleep(time.Duration(pi) * time.Second)
	}

}

func metricReporting(ri int64, r reporter) {
	for {
		err := r.ReportMetrics()
		if err != nil {
			fmt.Println(err)
		}
		time.Sleep(time.Duration(ri) * time.Second)
	}
}

func main() {
	err := parseFlags()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	ms := []string{
		"Alloc",
		"BuckHashSys",
		"Frees",
		"GCCPUFraction",
		"GCSys",
		"HeapAlloc",
		"HeapIdle",
		"HeapInuse",
		"HeapObjects",
		"HeapReleased",
		"HeapSys",
		"LastGC",
		"Lookups",
		"MCacheInuse",
		"MCacheSys",
		"MSpanInuse",
		"MSpanSys",
		"Mallocs",
		"NextGC",
		"NumForcedGC",
		"NumGC",
		"OtherSys",
		"PauseTotalNs",
		"StackInuse",
		"StackSys",
		"Sys",
		"TotalAlloc",
	}
	db := memstorage.NewMemStorage()
	collectorSvc := metrics.NewCollectorSvc(db)
	var ds bool
	if secretKey != "" {
		ds = true
		os.Setenv("CLIENT_SECRET_KEY", secretKey)
	} else {
		ds = false
	}
	reporter := http.NewReporter(serverAddress, db, ds)
	go metricCollecting(pollInterval, collectorSvc, ms)
	go metricReporting(reportInterval, reporter)
	fmt.Printf("Agent started\nReporting to: %s\nPollInterval: %d\nReportInterval: %d\n", serverAddress, pollInterval, reportInterval)
	select {}
}
