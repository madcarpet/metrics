package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/madcarpet/metrics/internal/adapter/memstorage"
	"github.com/madcarpet/metrics/internal/service/metrics"
	"github.com/stretchr/testify/assert"
)

func TestUpdateHandler(t *testing.T) {
	type want struct {
		code        int
		contentType string
		respBody    string
	}
	tests := []struct {
		name           string
		method         string
		reqContentType string
		reqBody        string
		url            string
		want           want
	}{
		{
			name:           "Test add new gauge",
			url:            "http://localhost:8080/update/",
			method:         http.MethodPost,
			reqContentType: "application/json",
			reqBody:        `{"id":"TestGauge1","type":"gauge","value":1212155}`,
			want: want{
				code:        http.StatusOK,
				contentType: "application/json",
				respBody:    `{"id":"TestGauge1","type":"gauge","value":1212155}`,
			},
		},
		{
			name:           "Test add new negtive gauge",
			url:            "http://localhost:8080/update/",
			method:         http.MethodPost,
			reqContentType: "application/json",
			reqBody:        `{"id":"TestGauge2","type":"gauge","value":-215111}`,
			want: want{
				code:        http.StatusOK,
				contentType: "application/json",
				respBody:    `{"id":"TestGauge2","type":"gauge","value":-215111}`,
			},
		},
		{
			name:           "Test add new gauge with point",
			url:            "http://localhost:8080/update/",
			method:         http.MethodPost,
			reqContentType: "application/json",
			reqBody:        `{"id":"TestGauge3","type":"gauge","value":0.12151244}`,
			want: want{
				code:        http.StatusOK,
				contentType: "application/json",
				respBody:    `{"id":"TestGauge3","type":"gauge","value":0.12151244}`,
			},
		},
		{
			name:           "Test add new negative gauge with point",
			url:            "http://localhost:8080/update/",
			method:         http.MethodPost,
			reqContentType: "application/json",
			reqBody:        `{"id":"TestGauge4","type":"gauge","value":-12.215111}`,
			want: want{
				code:        http.StatusOK,
				contentType: "application/json",
				respBody:    `{"id":"TestGauge4","type":"gauge","value":-12.215111}`,
			},
		},
		{
			name:           "Test update gauge",
			url:            "http://localhost:8080/update/",
			method:         http.MethodPost,
			reqContentType: "application/json",
			reqBody:        `{"id":"TestGauge1","type":"gauge","value":77.2152544}`,
			want: want{
				code:        http.StatusOK,
				contentType: "application/json",
				respBody:    `{"id":"TestGauge1","type":"gauge","value":77.2152544}`,
			},
		},
		{
			name:           "Test gauge without value",
			url:            "http://localhost:8080/update/",
			method:         http.MethodPost,
			reqContentType: "application/json",
			reqBody:        `{"id":"TestGauge1","type":"gauge"}`,
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain; charset=UTF-8",
				respBody:    "Bad request",
			},
		},
		{
			name:           "Test gauge with wrong value",
			url:            "http://localhost:8080/update/",
			method:         http.MethodPost,
			reqContentType: "application/json",
			reqBody:        `{"id":"TestGauge1","type":"gauge", "value":""}`,
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain; charset=UTF-8",
				respBody:    "Bad request",
			},
		},
		{
			name:           "Test add new counter",
			url:            "http://localhost:8080/update/",
			method:         http.MethodPost,
			reqContentType: "application/json",
			reqBody:        `{"id":"TestCounter1","type":"counter","delta":12}`,
			want: want{
				code:        http.StatusOK,
				contentType: "application/json",
				respBody:    `{"id":"TestCounter1","type":"counter","delta":12}`,
			},
		},
		{
			name:           "Test add new negative counter",
			url:            "http://localhost:8080/update/",
			method:         http.MethodPost,
			reqContentType: "application/json",
			reqBody:        `{"id":"TestCounter2","type":"counter","delta":50}`,
			want: want{
				code:        http.StatusOK,
				contentType: "application/json",
				respBody:    `{"id":"TestCounter2","type":"counter","delta":50}`,
			},
		},
		{
			name:           "Test add new zero counter",
			url:            "http://localhost:8080/update/",
			method:         http.MethodPost,
			reqContentType: "application/json",
			reqBody:        `{"id":"TestCounter3","type":"counter","delta":0}`,
			want: want{
				code:        http.StatusOK,
				contentType: "application/json",
				respBody:    `{"id":"TestCounter3","type":"counter","delta":0}`,
			},
		},
		{
			name:           "Test first update counter",
			url:            "http://localhost:8080/update/",
			method:         http.MethodPost,
			reqContentType: "application/json",
			reqBody:        `{"id":"TestCounter1","type":"counter","delta":10}`,
			want: want{
				code:        http.StatusOK,
				contentType: "application/json",
				respBody:    `{"id":"TestCounter1","type":"counter","delta":22}`,
			},
		},
		{
			name:           "Test second update counter",
			url:            "http://localhost:8080/update/",
			method:         http.MethodPost,
			reqContentType: "application/json",
			reqBody:        `{"id":"TestCounter1","type":"counter","delta":-2}`,
			want: want{
				code:        http.StatusOK,
				contentType: "application/json",
				respBody:    `{"id":"TestCounter1","type":"counter","delta":20}`,
			},
		},
		{
			name:           "Test counter without delta",
			url:            "http://localhost:8080/update/",
			method:         http.MethodPost,
			reqContentType: "application/json",
			reqBody:        `{"id":"TestCounter1","type":"counter"}`,
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain; charset=UTF-8",
				respBody:    "Bad request",
			},
		},
		{
			name:           "Test counter without delta and with value",
			url:            "http://localhost:8080/update/",
			method:         http.MethodPost,
			reqContentType: "application/json",
			reqBody:        `{"id":"TestCounter1","type":"counter","value":50}`,
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain; charset=UTF-8",
				respBody:    "Bad request",
			},
		},
		{
			name:           "Test counter wrong delta",
			url:            "http://localhost:8080/update/",
			method:         http.MethodPost,
			reqContentType: "application/json",
			reqBody:        `{"id":"TestCounter1","type":"counter","delta":""}`,
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain; charset=UTF-8",
				respBody:    "Bad request",
			},
		},
		{
			name:           "Test request without id",
			url:            "http://localhost:8080/update/",
			method:         http.MethodPost,
			reqContentType: "application/json",
			reqBody:        `{"type":"counter","delta":12}`,
			want: want{
				code:        http.StatusNotFound,
				contentType: "text/plain; charset=UTF-8",
				respBody:    "Metric name not found",
			},
		},
		{
			name:           "Test request with empty id",
			url:            "http://localhost:8080/update/",
			method:         http.MethodPost,
			reqContentType: "application/json",
			reqBody:        `{"id":"","type":"counter","delta":12}`,
			want: want{
				code:        http.StatusNotFound,
				contentType: "text/plain; charset=UTF-8",
				respBody:    "Metric name not found",
			},
		},
		{
			name:           "Test request without type",
			url:            "http://localhost:8080/update/",
			method:         http.MethodPost,
			reqContentType: "application/json",
			reqBody:        `{"id":"TestCounter1","delta":12}`,
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain; charset=UTF-8",
				respBody:    "Bad request",
			},
		},
		{
			name:           "Test wrong Content-Type",
			url:            "http://localhost:8080/update/",
			method:         http.MethodPost,
			reqContentType: "text/plain; charset=UTF-8",
			reqBody:        `{"id":"TestCounter1","type":"counter","delta":12}`,
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain; charset=UTF-8",
				respBody:    "Bad request",
			},
		},
	}
	db := memstorage.NewMemStorage()
	e := echo.New()
	updateSvc := metrics.NewUpdateMetricSvc(db)
	getSvc := metrics.NewGetMetricSvc(db)
	updateHandler := NewUpdateHandler(updateSvc, getSvc)
	e.POST("/update/", updateHandler.Handle)

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
			assert.Equal(t, test.want.contentType, rec.Header().Get("Content-Type"))
			resRespBody := strings.TrimRight(rec.Body.String(), "\n")
			assert.Equal(t, test.want.respBody, resRespBody)
		},
		)
	}
}
