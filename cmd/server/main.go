package main

import (
	"fmt"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/madcarpet/metrics/internal/adapter/memstorage"
	"github.com/madcarpet/metrics/internal/handlers/httpecho"
	"github.com/madcarpet/metrics/internal/logger"
	"github.com/madcarpet/metrics/internal/service/metrics"
)

func main() {
	err := parseFlags()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	logger.Initialize(loggingLevel)
	defer logger.Log.Sync()

	db := memstorage.NewMemStorage()
	rootSvc := metrics.NewGetAllMetricsSvc(db)
	valueSvc := metrics.NewGetMetricSvc(db)
	updateSvc := metrics.NewUpdateMetricSvc(db)

	e := echo.New()
	httpecho.SetupRouter(e, rootSvc, valueSvc, updateSvc)

	logger.Log.Info("Server starting")

	err = e.Start(serverAddress)
	if err != nil {
		logger.Log.Fatal("error starting server")
	}
}
