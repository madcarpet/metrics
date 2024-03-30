package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/madcarpet/metrics/internal/adapter/localdisk"
	"github.com/madcarpet/metrics/internal/adapter/memstorage"
	"github.com/madcarpet/metrics/internal/handlers/httpecho"
	"github.com/madcarpet/metrics/internal/logger"
	"github.com/madcarpet/metrics/internal/service/metrics"
)

type dbWriter interface {
	WriteDb() error
}

func writeDbToDisk(w dbWriter, i int64) {
	logger.Log.Info(fmt.Sprintf("Server data storing interval: %d\n", i))
	logger.Log.Info(fmt.Sprintf("Server will store data to file: %s", fileStoragePath))
	switch i {
	case 0:
		for {
			w.WriteDb()
		}
	default:
		for {
			w.WriteDb()
			time.Sleep(time.Duration(i) * time.Second)
		}
	}
}

func main() {
	err := parseFlags()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	//create channels for error and stopping
	sigChan := make(chan os.Signal, 1)
	errChan := make(chan error)
	//register system signals wuth channels
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	logger.Initialize(loggingLevel)
	defer logger.Log.Sync()

	db := memstorage.NewMemStorage()
	rootSvc := metrics.NewGetAllMetricsSvc(db)
	valueSvc := metrics.NewGetMetricSvc(db)
	updateSvc := metrics.NewUpdateMetricSvc(db)
	dbSaver := localdisk.NewSaver(fileStoragePath, db)

	if isRestore {
		err = localdisk.DBImport(fileStoragePath, updateSvc)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	if fileStoragePath != "" {
		go writeDbToDisk(dbSaver, storeInterval)
	}

	e := echo.New()
	httpecho.SetupRouter(e, rootSvc, valueSvc, updateSvc)
	logger.Log.Info("Server starting")

	// Start server as goroutine, and send error to channel
	go func() {
		errChan <- e.Start(serverAddress)
	}()

	//handle channels
	select {
	case stop := <-sigChan:
		fmt.Printf("Server stopping, recieved signal: %v\n", stop)
		if fileStoragePath != "" {
			dbSaver.WriteDb()
		}
	case err := <-errChan:
		if err != nil {
			fmt.Printf("Server starting error: %v\n", err)
			os.Exit(1)
		}

	}
}
