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

func ExampleValueURLHandler_Handle() {
	// Init storage.
	db := memstorage.NewMemStorage()
	testMetrics := []entity.Metric{
		{Name: "TestGauge1", Type: entity.Gauge, Value: 1.114112e+06},
		{Name: "TestCounter1", Type: entity.Counter, Value: 188},
	}

	for _, m := range testMetrics {
		db.UpdateMetric(context.Background(), m)
	}
	// Init services.
	valueSvc := metrics.NewGetMetricSvc(db)
	valueURLHandler := NewValueURLHandler(valueSvc)
	// Init echo framework.
	e := echo.New()
	// Set handler.
	e.GET("/value/:type/:name", valueURLHandler.Handle)
	// Prepare requests and recorders
	reqG := httptest.NewRequest(http.MethodGet, "http://localhost:8080/value/gauge/TestGauge1", nil)
	reqC := httptest.NewRequest(http.MethodGet, "http://localhost:8080/value/counter/TestCounter1", nil)
	recG := httptest.NewRecorder()
	recC := httptest.NewRecorder()
	defer recG.Result().Body.Close()
	defer recC.Result().Body.Close()
	// Handle requests.
	e.ServeHTTP(recG, reqG)
	e.ServeHTTP(recC, reqC)

	// Output response details
	fmt.Println("Status Code:", recG.Code)
	fmt.Println("Content Type:", recG.Header().Get("Content-Type"))
	fmt.Println("Response Body:")
	fmt.Println(recG.Body.String())

	fmt.Println("Status Code:", recC.Code)
	fmt.Println("Content Type:", recC.Header().Get("Content-Type"))
	fmt.Println("Response Body:")
	fmt.Println(recC.Body.String())

	// Output:
	// Status Code: 200
	// Content Type: text/plain; charset=UTF-8
	// Response Body:
	// 1.114112e+06
	// Status Code: 200
	// Content Type: text/plain; charset=UTF-8
	// Response Body:
	// 188

}
