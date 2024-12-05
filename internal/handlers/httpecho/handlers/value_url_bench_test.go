package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/madcarpet/metrics/internal/adapter/storage/memstorage"
	"github.com/madcarpet/metrics/internal/entity"
	"github.com/madcarpet/metrics/internal/service/metrics"
)

func BenchmarkXxx(b *testing.B) {
	db := memstorage.NewMemStorage()
	testMetrics := []entity.Metric{
		{Name: "TestGauge1", Type: entity.Gauge, Value: 1.114112e+06},
		{Name: "TestCounter1", Type: entity.Counter, Value: 188},
	}

	for _, m := range testMetrics {
		db.UpdateMetric(context.Background(), m)
	}

	valueSvc := metrics.NewGetMetricSvc(db)
	valueURLHandler := NewValueURLHandler(valueSvc)
	e := echo.New()
	e.GET("/value/:type/:name", valueURLHandler.Handle)

	reqG := httptest.NewRequest(http.MethodGet, "http://localhost:8080/value/gauge/TestGauge1", nil)
	reqC := httptest.NewRequest(http.MethodGet, "http://localhost:8080/value/counter/TestCounter1", nil)
	recG := httptest.NewRecorder()
	recC := httptest.NewRecorder()
	defer recG.Result().Body.Close()
	defer recC.Result().Body.Close()
	b.ResetTimer()
	b.Run("value url gauge benchmark", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			e.ServeHTTP(recG, reqG)
		}
	})
	b.Run("value url counter benchmark", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			e.ServeHTTP(recC, reqC)
		}
	})

}
