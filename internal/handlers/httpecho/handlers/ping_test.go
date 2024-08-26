package handlers

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/madcarpet/metrics/internal/mocks"
	"github.com/madcarpet/metrics/internal/service/metrics"
	"github.com/stretchr/testify/assert"
)

func TestPingHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	s := mocks.NewMockRepository(ctrl)
	s.EXPECT().IsConnected(context.Background()).Return(nil)

	e := echo.New()
	pingSvc := metrics.NewPingSvc(s)
	pingHandler := NewPingHandler(context.Background(), pingSvc)
	e.GET("/ping", pingHandler.Handle)

	req := httptest.NewRequest(http.MethodGet, "http://localhost:8080/ping", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	defer rec.Result().Body.Close()
	assert.Equal(t, http.StatusOK, rec.Code)

	s.EXPECT().IsConnected(context.Background()).Return(fmt.Errorf("Error"))
	req2 := httptest.NewRequest(http.MethodGet, "http://localhost:8080/ping", nil)
	rec2 := httptest.NewRecorder()
	e.ServeHTTP(rec2, req2)
	assert.Equal(t, http.StatusInternalServerError, rec2.Code)

}
