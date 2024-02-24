package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestAnyHandler(t *testing.T) {
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
			name:   "Test GET root",
			url:    "http://localhost:8080/12dsad",
			method: http.MethodGet,
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain; charset=UTF-8",
				body:        "Bad request",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			e := echo.New()
			e.Any("/*", func(c echo.Context) error {
				return c.String(http.StatusBadRequest, "Bad request")
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
