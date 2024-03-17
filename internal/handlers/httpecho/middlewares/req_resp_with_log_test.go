package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/madcarpet/metrics/internal/logger"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func TestReqRespWithLogging(t *testing.T) {
	observedZapCore, observedLogs := observer.New(zap.InfoLevel)

	logger.Log = zap.New(observedZapCore)

	e := echo.New()
	e.Any("/*", func(c echo.Context) error {
		c.Response().Header().Set("Content-Type", "text/plain; charset=UTF-8")
		return c.String(http.StatusBadRequest, "Bad request")
	}, ReqRespWithLogging)

	req := httptest.NewRequest(http.MethodGet, "http://localhost:8080/testing", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	defer rec.Result().Body.Close()

	assert.Equal(t, 2, observedLogs.Len())
	assert.Equal(t, "info", observedLogs.All()[0].Level.String())
	firstMessage := observedLogs.All()[0].Context
	secondMessage := observedLogs.All()[1].Context
	assert.Contains(t, firstMessage[0].String, "testing")
	assert.Contains(t, secondMessage[0].String, "400")

}
