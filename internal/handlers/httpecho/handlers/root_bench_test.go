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

func BenchmarkRootHandler(b *testing.B) {
	db := memstorage.NewMemStorage()
	testMetrics := []entity.Metric{
		{Name: "TestGauge1", Type: entity.Gauge, Value: 1.114112e+06},
		{Name: "TestGauge2", Type: entity.Gauge, Value: 879464},
		{Name: "TestGauge3", Type: entity.Gauge, Value: 0.8230922114274958},
		{Name: "TestCounter1", Type: entity.Counter, Value: 188},
		{Name: "TestCounter2", Type: entity.Counter, Value: 991112111111},
		{Name: "TestCounter3", Type: entity.Counter, Value: 0},
	}

	for _, m := range testMetrics {
		db.UpdateMetric(context.Background(), m)
	}

	e := echo.New()
	rootSvc := metrics.NewGetAllMetricsSvc(db)
	rootHandler := NewRootHandler(rootSvc)
	e.GET("/", rootHandler.Handle)
	req := httptest.NewRequest(http.MethodGet, "http://localhost:8080/", nil)
	rec := httptest.NewRecorder()
	defer rec.Result().Body.Close()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		e.ServeHTTP(rec, req)
	}
}
