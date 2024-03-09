package middlewares

import (
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/madcarpet/metrics/internal/logger"
	"go.uber.org/zap"
)

func ReqRespWithLogging(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()
		req := c.Request()
		err := next(c)
		duration := time.Since(start)
		resp := c.Response()
		logger.Log.Info("Request data", zap.String("URI", req.RequestURI), zap.String("Metod", req.Method), zap.Duration("Duration", time.Duration(duration)))
		logger.Log.Info("Response data", zap.String("Code", fmt.Sprintf("%d", resp.Status)), zap.String("Size", fmt.Sprintf("%d", resp.Size)))
		return err
	}
}
