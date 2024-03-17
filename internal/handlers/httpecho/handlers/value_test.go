package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/madcarpet/metrics/internal/adapter/memstorage"
	"github.com/madcarpet/metrics/internal/entity"
	"github.com/madcarpet/metrics/internal/service/metrics"
	"github.com/stretchr/testify/assert"
)

func TestValueHanler(t *testing.T) {
	type want struct {
		code            int
		respContentType string
		respBody        string
	}
	tests := []struct {
		name           string
		url            string
		method         string
		reqContentType string
		reqBody        string
		want           want
	}{
		{
			name:           "Test POST correct gauge 1114112",
			url:            "http://localhost:8080/value/",
			method:         http.MethodPost,
			reqContentType: "application/json",
			reqBody:        `{"id":"TestGauge1","type":"gauge"}`,
			want: want{
				code:            http.StatusOK,
				respContentType: "application/json",
				respBody:        `{"id":"TestGauge1","type":"gauge","value":1114112}`,
			},
		},
		{
			name:           "Test POST correct negative gauge -12544",
			url:            "http://localhost:8080/value/",
			method:         http.MethodPost,
			reqContentType: "application/json",
			reqBody:        `{"id":"TestGauge2","type":"gauge"}`,
			want: want{
				code:            http.StatusOK,
				respContentType: "application/json",
				respBody:        `{"id":"TestGauge2","type":"gauge","value":-12544}`,
			},
		},
		{
			name:           "Test POST correct negative gauge with point -0.2154442",
			url:            "http://localhost:8080/value/",
			method:         http.MethodPost,
			reqContentType: "application/json",
			reqBody:        `{"id":"TestGauge3","type":"gauge"}`,
			want: want{
				code:            http.StatusOK,
				respContentType: "application/json",
				respBody:        `{"id":"TestGauge3","type":"gauge","value":-0.2154442}`,
			},
		},
		{
			name:           "Test POST correct gauge with point 0.8230922114274958",
			url:            "http://localhost:8080/value/",
			method:         http.MethodPost,
			reqContentType: "application/json",
			reqBody:        `{"id":"TestGauge4","type":"gauge"}`,
			want: want{
				code:            http.StatusOK,
				respContentType: "application/json",
				respBody:        `{"id":"TestGauge4","type":"gauge","value":0.8230922114274958}`,
			},
		},
		{
			name:           "Test POST correct counter 188",
			url:            "http://localhost:8080/value/",
			method:         http.MethodPost,
			reqContentType: "application/json",
			reqBody:        `{"id":"TestCounter1","type":"counter"}`,
			want: want{
				code:            http.StatusOK,
				respContentType: "application/json",
				respBody:        `{"id":"TestCounter1","type":"counter","delta":188}`,
			},
		},
		{
			name:           "Test POST correct negative counter -50",
			url:            "http://localhost:8080/value/",
			method:         http.MethodPost,
			reqContentType: "application/json",
			reqBody:        `{"id":"TestCounter2","type":"counter"}`,
			want: want{
				code:            http.StatusOK,
				respContentType: "application/json",
				respBody:        `{"id":"TestCounter2","type":"counter","delta":-50}`,
			},
		},
		{
			name:           "Test POST correct negative counter zero",
			url:            "http://localhost:8080/value/",
			method:         http.MethodPost,
			reqContentType: "application/json",
			reqBody:        `{"id":"TestCounter3","type":"counter"}`,
			want: want{
				code:            http.StatusOK,
				respContentType: "application/json",
				respBody:        `{"id":"TestCounter3","type":"counter","delta":0}`,
			},
		},
		{
			name:           "Test wrong content-type",
			url:            "http://localhost:8080/value/",
			method:         http.MethodPost,
			reqContentType: "text/plain",
			reqBody:        `{"id":"TestGauge1","type":"gauge"}`,
			want: want{
				code:            http.StatusBadRequest,
				respContentType: "text/plain; charset=UTF-8",
				respBody:        "Bad request",
			},
		},
		{
			name:           "Test wrong metric name",
			url:            "http://localhost:8080/value/",
			method:         http.MethodPost,
			reqContentType: "application/json",
			reqBody:        `{"id":"TestGauge11","type":"gauge"}`,
			want: want{
				code:            http.StatusNotFound,
				respContentType: "text/plain; charset=UTF-8",
				respBody:        "Metric name not found",
			},
		},
		{
			name:           "Test empty metric name",
			url:            "http://localhost:8080/value/",
			method:         http.MethodPost,
			reqContentType: "application/json",
			reqBody:        `{"id":"","type":"gauge"}`,
			want: want{
				code:            http.StatusNotFound,
				respContentType: "text/plain; charset=UTF-8",
				respBody:        "Metric name not found",
			},
		},
		{
			name:           "Test request without id",
			url:            "http://localhost:8080/value/",
			method:         http.MethodPost,
			reqContentType: "application/json",
			reqBody:        `{"type":"gauge"}`,
			want: want{
				code:            http.StatusNotFound,
				respContentType: "text/plain; charset=UTF-8",
				respBody:        "Metric name not found",
			},
		},
		{
			name:           "Test request without type",
			url:            "http://localhost:8080/value/",
			method:         http.MethodPost,
			reqContentType: "application/json",
			reqBody:        `{"id":"TestGauge1"}`,
			want: want{
				code:            http.StatusBadRequest,
				respContentType: "text/plain; charset=UTF-8",
				respBody:        "Bad request",
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
			req.Header.Set("Content-Type", test.reqContentType)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)
			defer rec.Result().Body.Close()
			assert.Equal(t, test.want.code, rec.Code)
			assert.Equal(t, test.want.respContentType, rec.Header().Get("Content-Type"))
			resRespBody := strings.TrimRight(rec.Body.String(), "\n")
			assert.Equal(t, test.want.respBody, resRespBody)

		},
		)
	}
}
