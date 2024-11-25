package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/madcarpet/metrics/internal/mocks"
	"github.com/madcarpet/metrics/internal/service/metrics"
)

func BenchmarkPingHandler(b *testing.B) {
	ctrl := gomock.NewController(b)
	defer ctrl.Finish()
	s := mocks.NewMockRepository(ctrl)
	s.EXPECT().IsConnected(context.Background()).Return(nil).AnyTimes()

	e := echo.New()
	pingSvc := metrics.NewPingSvc(s)
	pingHandler := NewPingHandler(context.Background(), pingSvc)
	e.GET("/ping", pingHandler.Handle)

	req := httptest.NewRequest(http.MethodGet, "http://localhost:8080/ping", nil)
	rec := httptest.NewRecorder()
	defer rec.Result().Body.Close()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		e.ServeHTTP(rec, req)
	}

}
