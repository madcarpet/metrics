package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/madcarpet/metrics/internal/adapter/dispenser"
	"github.com/madcarpet/metrics/internal/adapter/http"
	"github.com/madcarpet/metrics/internal/adapter/storage/memstorage"
	"github.com/madcarpet/metrics/internal/entity"
	"github.com/madcarpet/metrics/internal/logger"
	"github.com/madcarpet/metrics/internal/service/metrics"
)

// Vars for ldflags.
var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
)

// reporter - interface with a method `ReportMetrics` to report metrics to http server.
type reporter interface {
	ReportMetrics(ctx context.Context, metrics []entity.Metric) error
}

// collectService - interface with a method `Collect` to collect metrics system metrics.
type collectService interface {
	Collect(ms []string) error
}

// metricCollecting - function that realize collecting metrics with selected interval.
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

// worker - function that realize worker to parallelize metric reporting process.
func worker(ctx context.Context, n int, rpt reporter, chIn <-chan []entity.Metric) {
	fmt.Printf("worker #%d started\n", n)
	for metric := range chIn {
		batchSize := 3
		for i := 0; i < len(metric); i += batchSize {
			end := i + batchSize
			if end > len(metric) {
				end = len(metric)
			}
			batch := metric[i:end]
			err := rpt.ReportMetrics(ctx, batch)
			if err != nil {
				fmt.Printf("worker report metrics error %s\n", err)
			}
		}
	}
	fmt.Printf("worker #%d finished\n", n)
}

func main() {
	// Print version information.
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)
	logger.Initialize("info")
	defer logger.Log.Sync()
	doneChan := make(chan struct{})
	// create channels for error and stopping.
	sigChan := make(chan os.Signal, 1)
	errChan := make(chan error)
	// register system signals with channels.
	signal.Notify(sigChan, os.Interrupt)

	agentConfig, err := parseFlags()
	if err != nil {
		fmt.Println(err)
		return
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
	if agentConfig.SecretKey != "" {
		ds = true
	} else {
		ds = false
	}
	reporter := http.NewReporter(agentConfig.ServerAddress, ds, agentConfig.SecretKey, agentConfig.PKey)
	disp := dispenser.NewMetricDispenser(db, comCh)
	go metricCollecting(doneChan, agentConfig.PollInterval, collectorSvc, ms)
	go metricCollecting(doneChan, agentConfig.PollInterval, perfCollectorSvc, mn)
	for i := 1; i < int(agentConfig.RateLimit)+1; i++ {
		go worker(context.Background(), i, reporter, comCh)
	}
	go disp.Dispense(context.Background(), agentConfig.ReportInterval)

	fmt.Printf("Agent started\nReporting to: %s\nPollInterval: %d\nReportInterval: %d\n", agentConfig.ServerAddress, agentConfig.PollInterval, agentConfig.ReportInterval)
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
