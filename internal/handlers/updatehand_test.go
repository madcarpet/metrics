package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/madcarpet/metrics/internal/storage"
	"github.com/stretchr/testify/assert"
)

func TestUpdateHandeler(t *testing.T) {
	type want struct {
		code        int
		contentType string
	}
	tests := []struct {
		name   string
		method string
		url    string
		want   want
	}{
		{
			name:   "Test correct gauge",
			url:    "http://localhost:8080/update/gauge/test/1.33",
			method: http.MethodPost,
			want: want{
				code:        200,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:   "Test correct counter",
			url:    "http://localhost:8080/update/counter/test/1",
			method: http.MethodPost,
			want: want{
				code:        200,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:   "Test incorrect gauge value",
			url:    "http://localhost:8080/update/gauge/test/b",
			method: http.MethodPost,
			want: want{
				code:        400,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:   "Test incorrect counter value",
			url:    "http://localhost:8080/update/counter/test/b",
			method: http.MethodPost,
			want: want{
				code:        400,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:   "Test missed metric name",
			url:    "http://localhost:8080/update/counter/12",
			method: http.MethodPost,
			want: want{
				code:        404,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:   "Test incorrect type",
			url:    "http://localhost:8080/update/bool/test/1",
			method: http.MethodPost,
			want: want{
				code:        400,
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := storage.NewMemStorage()
			req := httptest.NewRequest(test.method, test.url, nil)
			rec := httptest.NewRecorder()
			Update(rec, req, s)
			res := rec.Result()
			assert.Equal(t, test.want.code, res.StatusCode)
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))

		},
		)
	}
}
