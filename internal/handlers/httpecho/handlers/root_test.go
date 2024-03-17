package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/madcarpet/metrics/internal/adapter/memstorage"
	"github.com/madcarpet/metrics/internal/entity"
	"github.com/madcarpet/metrics/internal/service/metrics"
	"github.com/stretchr/testify/assert"
)

func TestRootHandler(t *testing.T) {
	type want struct {
		code        int
		contentType string
		body        []string
	}
	tests := []struct {
		name   string
		method string
		url    string
		want   want
	}{
		{
			name:   "Test root path",
			url:    "http://localhost:8080/",
			method: http.MethodGet,
			want: want{
				code:        http.StatusOK,
				contentType: "text/plain; charset=UTF-8",
				body:        []string{"1.114112e+06", "TestGauge2", "188", "TestCounter3"},
			},
		},
	}
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
		db.UpdateMetric(m)
	}

	e := echo.New()
	rootSvc := metrics.NewGetAllMetricsSvc(db)
	rootHandler := NewRootHandler(rootSvc)
	e.GET("/", rootHandler.Handle)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(test.method, test.url, nil)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)
			defer rec.Result().Body.Close()
			assert.Equal(t, test.want.code, rec.Code)
			assert.Equal(t, test.want.contentType, rec.Header().Get("Content-Type"))
			for _, k := range test.want.body {
				assert.Contains(t, rec.Body.String(), k)
			}

		},
		)
	}

}
