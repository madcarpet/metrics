package handlers

import (
	"net/http"
	"net/http/httptest"
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
		body        string
	}
	tests := []struct {
		name   string
		method string
		url    string
		want   want
	}{
		{
			name:   "Test update without name",
			url:    "http://localhost:8080/update/gauge//1.21211",
			method: http.MethodPost,
			want: want{
				code:        http.StatusNotFound,
				contentType: "text/plain; charset=UTF-8",
				body:        "Metric name not found",
			},
		},
		{
			name:   "Test wrong metric type",
			url:    "http://localhost:8080/update/bool/TestGauge1/1.21211",
			method: http.MethodPost,
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain; charset=UTF-8",
				body:        "Bad request",
			},
		},
		{
			name:   "Test gauge wrong value",
			url:    "http://localhost:8080/update/gauge/TestGauge1/badvalue",
			method: http.MethodPost,
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain; charset=UTF-8",
				body:        "Bad request",
			},
		},
		{
			name:   "Test gauge correct value",
			url:    "http://localhost:8080/update/gauge/TestGauge1/1.21511",
			method: http.MethodPost,
			want: want{
				code:        http.StatusOK,
				contentType: "text/plain; charset=UTF-8",
				body:        "Metric updated",
			},
		},
		{
			name:   "Test counter wrong value",
			url:    "http://localhost:8080/update/counter/TestCounter1/1.21511",
			method: http.MethodPost,
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain; charset=UTF-8",
				body:        "Bad request",
			},
		},
		{
			name:   "Test counter correct value",
			url:    "http://localhost:8080/update/counter/TestCounter1/188",
			method: http.MethodPost,
			want: want{
				code:        http.StatusOK,
				contentType: "text/plain; charset=UTF-8",
				body:        "Metric updated",
			},
		},
	}
	db := memstorage.NewMemStorage()
	e := echo.New()
	updateSvc := metrics.NewUpdateMetricSvc(db)
	updateHandler := NewUpdateHandler(updateSvc)
	e.POST("/update/:type/:name/:value", updateHandler.Handle)

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
