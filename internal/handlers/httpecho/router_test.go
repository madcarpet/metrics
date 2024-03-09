package httpecho

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestRouter(t *testing.T) {
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
			name:   "Test GET not handled path",
			url:    "http://localhost:8080/some/path",
			method: http.MethodGet,
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain; charset=UTF-8",
				body:        "Bad request",
			},
		},
		{
			name:   "Test POST not handled path",
			url:    "http://localhost:8080/some/path",
			method: http.MethodPost,
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain; charset=UTF-8",
				body:        "Bad request",
			},
		},
		{
			name:   "Test Update without metric name and value",
			url:    "http://localhost:8080/update/gauge/",
			method: http.MethodPost,
			want: want{
				code:        http.StatusNotFound,
				contentType: "text/plain; charset=UTF-8",
				body:        "Metric name not found",
			},
		},
		{
			name:   "Test Update without metric name",
			url:    "http://localhost:8080/update/gauge/1.21511",
			method: http.MethodPost,
			want: want{
				code:        http.StatusNotFound,
				contentType: "text/plain; charset=UTF-8",
				body:        "Metric name not found",
			},
		},
	}
	e := echo.New()
	e.POST("/update/:type/", func(c echo.Context) error {
		c.Response().Header().Set("Content-Type", "text/plain; charset=UTF-8")
		return c.String(http.StatusNotFound, "Metric name not found")
	})
	e.POST("/update/:type/:value", func(c echo.Context) error {
		c.Response().Header().Set("Content-Type", "text/plain; charset=UTF-8")
		return c.String(http.StatusNotFound, "Metric name not found")
	})
	e.Any("/*", func(c echo.Context) error {
		c.Response().Header().Set("Content-Type", "text/plain; charset=UTF-8")
		return c.String(http.StatusBadRequest, "Bad request")
	})

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
