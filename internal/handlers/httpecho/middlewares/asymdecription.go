package middlewares

import (
	"io"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/madcarpet/metrics/internal/constants"
	"github.com/madcarpet/metrics/internal/encryption/asymetric"
	"github.com/madcarpet/metrics/internal/logger"
	"go.uber.org/zap"
)

func AsymetricDecrypt(PKeyPath string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if PKeyPath == "" || c.Request().Method != echo.POST {
				err := next(c)
				return err
			}
			body, err := io.ReadAll(c.Request().Body)
			defer c.Request().Body.Close()
			if err != nil {
				logger.Log.Error("body reading error", zap.Error(err))
				c.Response().Header().Set("Content-Type", constants.ContentTypePlain)
				return c.String(http.StatusInternalServerError, "Server error - body reading error")
			}
			decBody, err := asymetric.DecryptWithPrivateKey(PKeyPath, body)
			if err != nil {
				logger.Log.Error("body decrypting error", zap.Error(err))
				c.Response().Header().Set("Content-Type", constants.ContentTypePlain)
				return c.String(http.StatusBadRequest, "Bad request - body decrypting error")
			}
			modifiedBodyReader := io.NopCloser(strings.NewReader(string(decBody)))
			c.Request().Body = modifiedBodyReader
			c.Request().ContentLength = int64(len(decBody))
			err = next(c)
			return err
		}
	}
}
