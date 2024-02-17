package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/madcarpet/metrics/internal/storage"
	"github.com/stretchr/testify/assert"
)

func TestValueHandler(t *testing.T) {
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
			url:    "http://localhost:8080/value/gauge/TestGCSys",
			method: http.MethodGet,
			want: want{
				code:        http.StatusOK,
				contentType: "text/plain; charset=UTF-8",
				body:        "1.114112e+06",
			},
		},
		{
			name:   "Test GET correct gauge 879464",
			url:    "http://localhost:8080/value/gauge/TestAlloc",
			method: http.MethodGet,
			want: want{
				code:        http.StatusOK,
				contentType: "text/plain; charset=UTF-8",
				body:        "879464",
			},
		},
		{
			name:   "Test GET correct gauge 0.8230922114274958",
			url:    "http://localhost:8080/value/gauge/TestRandomValue",
			method: http.MethodGet,
			want: want{
				code:        http.StatusOK,
				contentType: "text/plain; charset=UTF-8",
				body:        "0.8230922114274958",
			},
		},
		{
			name:   "Test correct counter value",
			url:    "http://localhost:8080/value/counter/TestPollCount",
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
	db := storage.MemStorage{
		Gauges: map[string]float64{
			"TestGCSys":       1.114112e+06,
			"TestAlloc":       879464,
			"TestRandomValue": 0.8230922114274958,
		},
		Counters: map[string]int64{
			"TestPollCount": 188,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			e := echo.New()
			e.GET("/value/:type/:name", func(c echo.Context) error {
				return Value(c, &db)
			})
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
