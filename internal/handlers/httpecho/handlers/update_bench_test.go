package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/madcarpet/metrics/internal/adapter/storage/memstorage"
	"github.com/madcarpet/metrics/internal/service/metrics"
)

func BenchmarkUpdateHandler(b *testing.B) {
	db := memstorage.NewMemStorage()
	e := echo.New()
	updateSvc := metrics.NewUpdateMetricSvc(db)
	getSvc := metrics.NewGetMetricSvc(db)
	updateHandler := NewUpdateHandler(updateSvc, getSvc)
	e.POST("/update/", updateHandler.Handle)
	g := `{"id":"TestGauge1","type":"gauge","value":1.1212155}`
	c := `{"id":"TestCounter1","type":"counter","delta":12}`
	var bodyGBuffer bytes.Buffer
	bodyGBuffer.Write([]byte(g))
	var bodyCBuffer bytes.Buffer
	bodyCBuffer.Write([]byte(c))
	reqG := httptest.NewRequest(http.MethodPost, "http://localhost:8080/update/", &bodyGBuffer)
	reqC := httptest.NewRequest(http.MethodPost, "http://localhost:8080/update/", &bodyCBuffer)
	reqG.Header.Set("Content-Type", "application/json")
	reqC.Header.Set("Content-Type", "application/json")
	recG := httptest.NewRecorder()
	recC := httptest.NewRecorder()
	defer recG.Result().Body.Close()
	defer recC.Result().Body.Close()
	b.ResetTimer()
	b.Run("update gauge benchmark", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			e.ServeHTTP(recG, reqG)
		}
	})
	b.Run("update counter benchmark", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			e.ServeHTTP(recC, reqC)
		}
	})

}
