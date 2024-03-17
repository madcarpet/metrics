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

func TestValueUrlHanler(t *testing.T) {
	type want struct {
		code        int
		contentType string
		body        string
	}
	tests := []struct {
		name   string
		method string
		url    string
		want   want
	}{
		{
			name:   "Test GET correct gauge 1.114112e+06",
			url:    "http://localhost:8080/value/gauge/TestGauge1",
			method: http.MethodGet,
			want: want{
				code:        http.StatusOK,
				contentType: "text/plain; charset=UTF-8",
				body:        "1.114112e+06",
			},
		},
		{
			name:   "Test GET correct gauge 879464",
			url:    "http://localhost:8080/value/gauge/TestGauge2",
			method: http.MethodGet,
			want: want{
				code:        http.StatusOK,
				contentType: "text/plain; charset=UTF-8",
				body:        "879464",
			},
		},
		{
			name:   "Test GET correct gauge 0.8230922114274958",
			url:    "http://localhost:8080/value/gauge/TestGauge3",
			method: http.MethodGet,
			want: want{
				code:        http.StatusOK,
				contentType: "text/plain; charset=UTF-8",
				body:        "0.8230922114274958",
			},
		},
		{
			name:   "Test correct counter value",
			url:    "http://localhost:8080/value/counter/TestCounter1",
			method: http.MethodGet,
			want: want{
				code:        http.StatusOK,
				contentType: "text/plain; charset=UTF-8",
				body:        "188",
			},
		},
		{
			name:   "Test incorrect gauge name",
			url:    "http://localhost:8080/value/gauge/wefdfg",
			method: http.MethodGet,
			want: want{
				code:        http.StatusNotFound,
				contentType: "text/plain; charset=UTF-8",
				body:        "Metric name not found",
			},
		},
		{
			name:   "Test incorrect counter name",
			url:    "http://localhost:8080/value/counter/wefdfg",
			method: http.MethodGet,
			want: want{
				code:        http.StatusNotFound,
				contentType: "text/plain; charset=UTF-8",
				body:        "Metric name not found",
			},
		},
		{
			name:   "Test incorrect metric type",
			url:    "http://localhost:8080/value/sdsd/TestPollCount",
			method: http.MethodGet,
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain; charset=UTF-8",
				body:        "Bad request",
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
	valueSvc := metrics.NewGetMetricSvc(db)
	valueUrlHandler := NewValueUrlHandler(valueSvc)
	e.GET("/value/:type/:name", valueUrlHandler.Handle)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(test.method, test.url, nil)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)
			defer rec.Result().Body.Close()
			assert.Equal(t, test.want.code, rec.Code)
			assert.Equal(t, test.want.contentType, rec.Header().Get("Content-Type"))
			assert.Equal(t, test.want.body, rec.Body.String())

		},
		)
	}
}
