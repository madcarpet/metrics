package http

import (
	"net/http"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/madcarpet/metrics/internal/adapter/memstorage"
	"github.com/madcarpet/metrics/internal/entity"
	"github.com/stretchr/testify/assert"
)

func TestReporter(t *testing.T) {
	db := memstorage.NewMemStorage()
	testMetrics := []entity.Metric{
		{Name: "TestGauge1", Type: entity.Gauge, Value: 1.114112e+06},
		{Name: "TestCounter1", Type: entity.Counter, Value: 188},
	}
	for _, metric := range testMetrics {
		db.UpdateMetric(metric)
	}
	e := echo.New()
	anyHandler := func(c echo.Context) error {
		c.Response().Header().Set("Content-Type", "text/plain; charset=UTF-8")
		return c.String(http.StatusOK, "Metric updated")
	}
	e.Any("/*", anyHandler)
	go e.Start("localhost:8080")

	reporter := NewReporter("localhost:8080", db)
	err := reporter.ReportMetrics()
	assert.Nil(t, err)

	time.Sleep(5 * time.Second)
	e.Close()
}
