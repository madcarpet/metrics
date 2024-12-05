package handlers

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/labstack/echo/v4"
	"github.com/madcarpet/metrics/internal/adapter/storage/memstorage"
	"github.com/madcarpet/metrics/internal/service/metrics"
)

func ExampleUpdatesHandler_Handle() {
	// Init DB.
	db := memstorage.NewMemStorage()
	// Prepare request body.
	reqBody := `[{"id":"testUpdatesGA1","type":"gauge","value":500},{"id":"testUpdatesGA1","type":"gauge","value":1000},{"id":"testUpdatesGA1","type":"gauge","value":3000},{"id":"testUpdatesCO1","type":"counter","delta":10},{"id":"testUpdatesCO1","type":"counter","delta":15},{"id":"testUpdatesCO1","type":"counter","delta":25}]`
	// Init echo framework.
	e := echo.New()
	// Init services.
	updatesSvc := metrics.NewUpdateMetricsSvc(db)
	getSvc := metrics.NewGetMetricSvc(db)
	updatesHandler := NewUpdatesHandler(updatesSvc, getSvc)
	// Set handler.
	e.POST("/updates/", updatesHandler.Handle)
	// Prepare requests and recorders.
	var bodyBuffer bytes.Buffer
	bodyBuffer.Write([]byte(reqBody))
	req := httptest.NewRequest(http.MethodPost, "http://localhost:8080/updates/", &bodyBuffer)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	// Handle request.
	e.ServeHTTP(rec, req)

	// Output response details
	fmt.Println("Status Code:", rec.Code)
	fmt.Println("Content Type:", rec.Header().Get("Content-Type"))
	fmt.Println("Response Body:")
	fmt.Println(rec.Body.String())

	// Output:
	// Status Code: 200
	// Content Type: application/json
	// Response Body:
	// [{"id":"testUpdatesGA1","type":"gauge","value":3000},{"id":"testUpdatesCO1","type":"counter","delta":50}]

}
