package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

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

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{}"))
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()
	reporter := NewReporter(server.URL[7:], db)
	err := reporter.ReportMetrics()
	assert.Nil(t, err)

}
