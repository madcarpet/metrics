package handlers

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/madcarpet/metrics/internal/adapter/storage/memstorage"
	"github.com/madcarpet/metrics/internal/entity"
	"github.com/madcarpet/metrics/internal/service/metrics"
)

func BenchmarkValueHandler(b *testing.B) {
	db := memstorage.NewMemStorage()
	testMetrics := []entity.Metric{
		{Name: "TestGauge1", Type: entity.Gauge, Value: 1.114112e+06},
		{Name: "TestCounter1", Type: entity.Counter, Value: 188},
	}

	for _, m := range testMetrics {
		db.UpdateMetric(context.Background(), m)
	}
	e := echo.New()
	valueSvc := metrics.NewGetMetricSvc(db)
	valueHandler := NewValueHandler(valueSvc)
	e.POST("/value/", valueHandler.Handle)
	g := `{"id":"TestGauge1","type":"gauge"}`
	c := `{"id":"TestCounter1","type":"counter"}`
	var bodyGBuffer bytes.Buffer
	bodyGBuffer.Write([]byte(g))
	var bodyCBuffer bytes.Buffer
	bodyCBuffer.Write([]byte(c))
	reqG := httptest.NewRequest(http.MethodPost, "http://localhost:8080/value/", &bodyGBuffer)
	reqC := httptest.NewRequest(http.MethodPost, "http://localhost:8080/value/", &bodyCBuffer)
	recG := httptest.NewRecorder()
	recC := httptest.NewRecorder()
	defer recG.Result().Body.Close()
	defer recC.Result().Body.Close()
	b.ResetTimer()
	b.Run("value gauge benchmark", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			e.ServeHTTP(recG, reqG)
		}
	})
	b.Run("value counter benchmark", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			e.ServeHTTP(recC, reqC)
		}
	})
}
