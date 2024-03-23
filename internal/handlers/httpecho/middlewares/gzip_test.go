package middlewares

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/madcarpet/metrics/internal/adapter/memstorage"
	"github.com/madcarpet/metrics/internal/handlers/httpecho/handlers"
	"github.com/madcarpet/metrics/internal/service/metrics"
	"github.com/stretchr/testify/assert"
)

func TestGzipMiddleware(t *testing.T) {
	data := `{"id": "test2", "type": "gauge", "value": 15.212}`
	respData := `{"id":"test2","type":"gauge","value":15.212}`
	var bodyBuffer bytes.Buffer
	var respBody []byte
	gzWriter := gzip.NewWriter(&bodyBuffer)
	gzWriter.Write([]byte(data))
	gzWriter.Close()

	db := memstorage.NewMemStorage()
	e := echo.New()
	updateSvc := metrics.NewUpdateMetricSvc(db)
	getSvc := metrics.NewGetMetricSvc(db)
	updateHandler := handlers.NewUpdateHandler(updateSvc, getSvc)
	e.POST("/update/", updateHandler.Handle, GzipCompression)

	req := httptest.NewRequest(http.MethodPost, "http://localhost:8080/update/", &bodyBuffer)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept-Encoding", "gzip")
	req.Header.Set("Content-Encoding", "gzip")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	resp := rec.Result()
	resp.Body.Close()

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "gzip", rec.Header().Get("Content-Encoding"))

	gzReader, _ := gzip.NewReader(resp.Body)
	defer gzReader.Close()
	respBody, _ = io.ReadAll(gzReader)
	//testing response body after decompressing
	assert.Equal(t, respData, strings.TrimRight(string(respBody), "\n"))

}
