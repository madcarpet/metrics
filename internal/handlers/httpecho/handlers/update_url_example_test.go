package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/labstack/echo/v4"
	"github.com/madcarpet/metrics/internal/adapter/storage/memstorage"
	"github.com/madcarpet/metrics/internal/service/metrics"
)

func ExampleUpdateURLHandler_Handle() {
	// Init storage and echo framework.
	db := memstorage.NewMemStorage()
	e := echo.New()
	// Init services.
	updateSvc := metrics.NewUpdateMetricSvc(db)
	updateURLHandler := NewUpdateURLHandler(updateSvc)
	// Set handler.
	e.POST("/update/:type/:name/:value", updateURLHandler.Handle)
	// Prepare requests and recorders.
	reqG := httptest.NewRequest(http.MethodPost, "http://localhost:8080/update/gauge/TestGauge1/1.21511", nil)
	reqC := httptest.NewRequest(http.MethodPost, "http://localhost:8080/update/counter/TestCounter1/188", nil)
	recG := httptest.NewRecorder()
	recC := httptest.NewRecorder()
	defer recG.Result().Body.Close()
	defer recC.Result().Body.Close()
	// Handle request.
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
	// Metric updated
	// Status Code: 200
	// Content Type: text/plain; charset=UTF-8
	// Response Body:
	// Metric updated
}
