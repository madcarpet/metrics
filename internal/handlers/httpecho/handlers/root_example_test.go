package handlers

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/labstack/echo/v4"
	"github.com/madcarpet/metrics/internal/adapter/storage/memstorage"
	"github.com/madcarpet/metrics/internal/entity"
	"github.com/madcarpet/metrics/internal/service/metrics"
)

func ExampleRootHandler_Handle() {
	// Init storage and echo framework.
	db := memstorage.NewMemStorage()
	testMetrics := []entity.Metric{
		{Name: "TestGauge1", Type: entity.Gauge, Value: 1.114112e+06},
		{Name: "TestCounter1", Type: entity.Counter, Value: 188},
	}
	for _, m := range testMetrics {
		db.UpdateMetric(context.Background(), m)
	}
	e := echo.New()
	// Init services.
	rootSvc := metrics.NewGetAllMetricsSvc(db)
	rootHandler := NewRootHandler(rootSvc)
	// Set Handler.
	e.GET("/", rootHandler.Handle)
	// Create request and recorder.
	req := httptest.NewRequest(http.MethodGet, "http://localhost:8080/", nil)
	rec := httptest.NewRecorder()
	defer rec.Result().Body.Close()
	// Serve request.
	e.ServeHTTP(rec, req)

	// Output response details
	fmt.Println("Status Code:", rec.Code)
	fmt.Println("Content Type:", rec.Header().Get("Content-Type"))
	fmt.Println("Response Body:")
	fmt.Println(rec.Body.String())

	// Output:
	// Status Code: 200
	// Content Type: text/html
	// Response Body:
	// TestGauge1: 1.114112e+06
	// TestCounter1: 188

}
