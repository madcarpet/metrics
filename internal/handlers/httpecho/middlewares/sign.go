package middlewares

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

func signatory(data []byte, secretKey string) string {
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write(data)
	sign := h.Sum(nil)
	return hex.EncodeToString(sign)
}

type signedWriter struct {
	http.ResponseWriter
	writer *bytes.Buffer
}

func (sw *signedWriter) Write(d []byte) (int, error) {
	return sw.writer.Write(d)
}

func (sw *signedWriter) WriteHeader(statusCode int) {
}

func SignData(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		key := os.Getenv("SECRET_KEY")
		reqHeadSign := c.Request().Header.Values("HashSHA256")
		if len(reqHeadSign) == 0 || len(reqHeadSign) > 1 {
			return c.String(http.StatusBadRequest, "Bad request, no sign")
		} else {
			body, err := io.ReadAll(c.Request().Body)
			defer c.Request().Body.Close()
			if err != nil {
				c.Response().Header().Set("Content-Type", "text/plain; charset=UTF-8")
				return c.String(http.StatusInternalServerError, "Server error")
			}
			reqCalcSign := signatory(body, key)
			if reqHeadSign[0] != reqCalcSign {
				return c.String(http.StatusBadRequest, "Bad request, invalid sign")
			}
			c.Request().Body = io.NopCloser(bytes.NewBuffer(body))

		}
		var respBody bytes.Buffer
		sw := signedWriter{ResponseWriter: c.Response().Writer, writer: &respBody}
		c.Response().Writer = &sw
		err := next(c)
		if err != nil {
			return err
		}
		sign := signatory(respBody.Bytes(), key)
		c.Response().Writer = sw.ResponseWriter
		c.Response().Header().Set("HashSHA256", sign)
		_, err = c.Response().Writer.Write(respBody.Bytes())
		return err
	}
}
