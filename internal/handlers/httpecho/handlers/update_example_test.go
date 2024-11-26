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

func ExampleUpdateHandler_Handle() {
	// Init storage and echo framework.
	db := memstorage.NewMemStorage()
	e := echo.New()
	// Init services.
	updateSvc := metrics.NewUpdateMetricSvc(db)
	getSvc := metrics.NewGetMetricSvc(db)
	updateHandler := NewUpdateHandler(updateSvc, getSvc)
	// Set Handler.
	e.POST("/update/", updateHandler.Handle)
	// Set test data.
	g := `{"id":"TestGauge1","type":"gauge","value":1.1212155}`
	c := `{"id":"TestCounter1","type":"counter","delta":12}`
	var bodyGBuffer bytes.Buffer
	bodyGBuffer.Write([]byte(g))
	var bodyCBuffer bytes.Buffer
	bodyCBuffer.Write([]byte(c))
	// Set requests and recorders.
	reqG := httptest.NewRequest(http.MethodPost, "http://localhost:8080/update/", &bodyGBuffer)
	reqC := httptest.NewRequest(http.MethodPost, "http://localhost:8080/update/", &bodyCBuffer)
	reqG.Header.Set("Content-Type", "application/json")
	reqC.Header.Set("Content-Type", "application/json")
	recG := httptest.NewRecorder()
	recC := httptest.NewRecorder()
	defer recG.Result().Body.Close()
	defer recC.Result().Body.Close()
	//Handle requests
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
	// {"id":"TestGauge1","type":"gauge","value":1.1212155}
	//
	// Status Code: 200
	// Content Type: application/json
	// Response Body:
	// {"id":"TestCounter1","type":"counter","delta":12}
}
