package handlers

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/labstack/echo/v4"
	"github.com/madcarpet/metrics/internal/adapter/storage/memstorage"
	"github.com/madcarpet/metrics/internal/entity"
	"github.com/madcarpet/metrics/internal/service/metrics"
)

func ExampleValueHandler_Handle() {
	// Init storage.
	db := memstorage.NewMemStorage()
	testMetrics := []entity.Metric{
		{Name: "TestGauge1", Type: entity.Gauge, Value: 1.114112e+06},
		{Name: "TestCounter1", Type: entity.Counter, Value: 188},
	}

	for _, m := range testMetrics {
		db.UpdateMetric(context.Background(), m)
	}
	// Init echo framework.
	e := echo.New()
	// Init services.
	valueSvc := metrics.NewGetMetricSvc(db)
	valueHandler := NewValueHandler(valueSvc)
	// Set handler.
	e.POST("/value/", valueHandler.Handle)
	// Prepare requests and recorders
	g := `{"id":"TestGauge1","type":"gauge"}`
	c := `{"id":"TestCounter1","type":"counter"}`
	var bodyGBuffer bytes.Buffer
	bodyGBuffer.Write([]byte(g))
	var bodyCBuffer bytes.Buffer
	bodyCBuffer.Write([]byte(c))
	reqG := httptest.NewRequest(http.MethodPost, "http://localhost:8080/value/", &bodyGBuffer)
	reqC := httptest.NewRequest(http.MethodPost, "http://localhost:8080/value/", &bodyCBuffer)
	reqG.Header.Set("Content-Type", "application/json")
	reqC.Header.Set("Content-Type", "application/json")
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
	// Content Type: application/json
	// Response Body:
	// {"id":"TestGauge1","type":"gauge","value":1114112}
	//
	// Status Code: 200
	// Content Type: application/json
	// Response Body:
	// {"id":"TestCounter1","type":"counter","delta":188}
}
