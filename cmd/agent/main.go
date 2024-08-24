package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/madcarpet/metrics/internal/adapter/dispenser"
	"github.com/madcarpet/metrics/internal/adapter/http"
	"github.com/madcarpet/metrics/internal/adapter/storage/memstorage"
	"github.com/madcarpet/metrics/internal/entity"
	"github.com/madcarpet/metrics/internal/service/metrics"
)

type reporter interface {
	ReportMetrics(metrics []entity.Metric) error
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

func worker(n int, rpt reporter, chIn <-chan []entity.Metric) {
	fmt.Printf("worker #%d started\n", n)
	for metric := range chIn {
		rpt.ReportMetrics(metric)
	}
	fmt.Printf("worker #%d finished\n", n)
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
	mn := []string{
		"TotalMemory",
		"FreeMemory",
		"CPUutilization",
	}
	comCh := make(chan []entity.Metric)
	db := memstorage.NewMemStorage()
	collectorSvc := metrics.NewCollectorSvc(db)
	perfCollectorSvc := metrics.NewPerCollectorSvc(db)
	var ds bool
	if secretKey != "" {
		ds = true
		os.Setenv("CLIENT_SECRET_KEY", secretKey)
	} else {
		ds = false
	}
	reporter := http.NewReporter(serverAddress, ds)
	disp := dispenser.NewMetricDispenser(db, comCh)
	go metricCollecting(pollInterval, collectorSvc, ms)
	go metricCollecting(pollInterval, perfCollectorSvc, mn)
	for i := 1; i < int(rateLimit)+1; i++ {
		go worker(i, reporter, comCh)
	}
	go disp.Dispense(context.Background(), reportInterval)

	fmt.Printf("Agent started\nReporting to: %s\nPollInterval: %d\nReportInterval: %d\n", serverAddress, pollInterval, reportInterval)
	select {}
}
