package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/madcarpet/metrics/internal/storage"
	"github.com/stretchr/testify/assert"
)

func TestRootHandler(t *testing.T) {
	type want struct {
		code        int
		contentType string
		body        map[string]string
	}
	tests := []struct {
		name   string
		method string
		url    string
		want   want
	}{
		{
			name:   "Test GET root",
			url:    "http://localhost:8080/",
			method: http.MethodGet,
			want: want{
				code:        http.StatusOK,
				contentType: "text/plain; charset=UTF-8",
				body: map[string]string{
					"TestGCSys":       "1.114112e+06",
					"TestAlloc":       "879464",
					"TestRandomValue": "0.8230922114274958",
				},
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
			e.GET("/", func(c echo.Context) error {
				return Root(c, &db)
			})
			req := httptest.NewRequest(test.method, test.url, nil)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)
			defer rec.Result().Body.Close()
			assert.Equal(t, test.want.code, rec.Code)
			assert.Equal(t, test.want.contentType, rec.Header().Get("Content-Type"))
			for k, v := range test.want.body {
				assert.Contains(t, rec.Body.String(), k+": "+v)
			}
		},
		)
	}
}
