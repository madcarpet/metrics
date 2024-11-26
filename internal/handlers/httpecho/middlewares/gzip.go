package middlewares

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/madcarpet/metrics/internal/constants"
)

// gzWriter - new writer with custom Writer(with gzip) for switch in c.Response().Writer
type gzWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

// Write - custom Write method with compression.
func (w gzWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

// gzReader new reader with custom Reader(with gzip) for switch in c.Request().Body.
type gzReader struct {
	io.ReadCloser
	Reader io.Reader
}

// Read custom method with decompression.
func (r gzReader) Read(p []byte) (n int, err error) {
	return r.Reader.Read(p)
}

// GzipCompression middleware function to compress data.
func GzipCompression(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		aEncod := c.Request().Header.Values("Accept-Encoding")
		if len(aEncod) > 0 {
			for _, h := range aEncod {
				if strings.Contains(h, "gzip") {
					// gzip Writer.
					gz := gzip.NewWriter(c.Response().Writer)
					defer gz.Close()
					c.Response().Header().Set("Content-Encoding", "gzip")
					// switch c.Response().Writer to custom gzWriter with gzip compression.
					c.Response().Writer = gzWriter{ResponseWriter: c.Response().Writer, Writer: gz}
				}
			}
		}
		cEncod := c.Request().Header.Get("Content-Encoding")
		if strings.Contains(cEncod, "gzip") {
			// gzip Reader.
			rgz, err := gzip.NewReader(c.Request().Body)
			if err != nil {
				c.Response().Header().Set("Content-Type", constants.ContentTypePlain)
				return c.String(http.StatusInternalServerError, "Server error")
			}
			defer rgz.Close()
			// switch c.Request().Body to custom gzReader with decompression.
			c.Request().Body = gzReader{ReadCloser: c.Request().Body, Reader: rgz}
		}
		return next(c)
	}
}
