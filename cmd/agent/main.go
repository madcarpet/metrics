package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/madcarpet/metrics/internal/adapter/dispenser"
	"github.com/madcarpet/metrics/internal/adapter/http"
	"github.com/madcarpet/metrics/internal/adapter/storage/memstorage"
	"github.com/madcarpet/metrics/internal/entity"
	"github.com/madcarpet/metrics/internal/service/metrics"
)

type reporter interface {
	ReportMetrics(ctx context.Context, metrics []entity.Metric) error
}

type collectService interface {
	Collect(ms []string) error
}

func metricCollecting(dch <-chan struct{}, pi int64, c collectService, ms []string) {
	tick := time.NewTicker(time.Duration(pi) * time.Second)
	for {
		select {
		case <-dch:
			return
		case <-tick.C:
			c.Collect(ms)
		}
	}

}

func worker(ctx context.Context, n int, rpt reporter, chIn <-chan []entity.Metric) {
	fmt.Printf("worker #%d started\n", n)
	for metric := range chIn {
		rpt.ReportMetrics(ctx, metric)
	}
	fmt.Printf("worker #%d finished\n", n)
}

func main() {
	doneChan := make(chan struct{})
	//create channels for error and stopping
	sigChan := make(chan os.Signal, 1)
	errChan := make(chan error)
	//register system signals with channels
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

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
	} else {
		ds = false
	}
	reporter := http.NewReporter(serverAddress, ds, secretKey)
	disp := dispenser.NewMetricDispenser(db, comCh)
	go metricCollecting(doneChan, pollInterval, collectorSvc, ms)
	go metricCollecting(doneChan, pollInterval, perfCollectorSvc, mn)
	for i := 1; i < int(rateLimit)+1; i++ {
		go worker(context.Background(), i, reporter, comCh)
	}
	go disp.Dispense(context.Background(), reportInterval)

	fmt.Printf("Agent started\nReporting to: %s\nPollInterval: %d\nReportInterval: %d\n", serverAddress, pollInterval, reportInterval)
	select {
	case stop := <-sigChan:
		fmt.Printf("Server stopping, recieved signal: %v\n", stop)
		close(doneChan)
	case err := <-errChan:
		if err != nil {
			fmt.Printf("Server got error: %v\n", err)
			close(doneChan)
		}

	}
}
