package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/madcarpet/metrics/internal/adapter/memstorage"
	"github.com/madcarpet/metrics/internal/entity"
	"github.com/madcarpet/metrics/internal/service/metrics"
	"github.com/stretchr/testify/assert"
)

func TestValueHanler(t *testing.T) {
	type want struct {
		code        int
		contentType string
		respBody    string
	}
	tests := []struct {
		name    string
		url     string
		method  string
		reqBody string
		want    want
	}{
		{
			name:    "Test GET correct gauge 1114112",
			url:     "http://localhost:8080/value/",
			method:  http.MethodPost,
			reqBody: `{"id":"TestGauge1","type":"gauge"}`,
			want: want{
				code:        http.StatusOK,
				contentType: "application/json",
				respBody:    `{"id":"TestGauge1","type":"gauge","value":1114112}`,
			},
		},
	}
	db := memstorage.NewMemStorage()
	testMetrics := []entity.Metric{
		{Name: "TestGauge1", Type: entity.Gauge, Value: 1114112},
		{Name: "TestGauge2", Type: entity.Gauge, Value: -12544},
		{Name: "TestGauge3", Type: entity.Gauge, Value: -0.2154442},
		{Name: "TestGauge4", Type: entity.Gauge, Value: 0.8230922114274958},
		{Name: "TestCounter1", Type: entity.Counter, Value: 188},
		{Name: "TestCounter2", Type: entity.Counter, Value: -50},
		{Name: "TestCounter3", Type: entity.Counter, Value: 0},
	}

	for _, m := range testMetrics {
		db.UpdateMetric(m)
	}

	e := echo.New()
	valueSvc := metrics.NewGetMetricSvc(db)
	valueHandler := NewValueHandler(valueSvc)
	e.POST("/value/", valueHandler.Handle)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var bodyBuffer bytes.Buffer
			bodyBuffer.Write([]byte(test.reqBody))
			req := httptest.NewRequest(test.method, test.url, &bodyBuffer)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)
			defer rec.Result().Body.Close()
			assert.Equal(t, test.want.code, rec.Code)
			assert.Equal(t, test.want.contentType, rec.Header().Get("Content-Type"))
			assert.Equal(t, test.want.respBody, rec.Body.String())

		},
		)
	}
}
