package main

import (
	"fmt"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/madcarpet/metrics/internal/adapter/memstorage"
	"github.com/madcarpet/metrics/internal/handlers/httpecho"
	"github.com/madcarpet/metrics/internal/service/metrics"
)

func main() {
	err := parseFlags()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	db := memstorage.NewMemStorage()
	rootSvc := metrics.NewGetAllMetricsSvc(db)
	valueSvc := metrics.NewGetMetricSvc(db)
	updateSvc := metrics.NewUpdateMetricSvc(db)
	e := echo.New()
	httpecho.SetupRouter(e, rootSvc, valueSvc, updateSvc)
	e.Logger.Fatal(e.Start(serverAddress))
}
