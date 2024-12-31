package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/madcarpet/metrics/internal/entity"
	"github.com/stretchr/testify/assert"
)

func TestReporter(t *testing.T) {
	testMetrics := []entity.Metric{
		{Name: "TestGauge1", Type: entity.Gauge, Value: 1.114112e+06},
		{Name: "TestCounter1", Type: entity.Counter, Value: 188},
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{}"))
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()
	reporter := NewReporter(server.URL[7:], false, "", "")
	err := reporter.ReportMetrics(context.Background(), testMetrics)
	assert.Nil(t, err)

}
