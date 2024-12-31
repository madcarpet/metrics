package middlewares

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/madcarpet/metrics/internal/adapter/storage/memstorage"
	"github.com/madcarpet/metrics/internal/handlers/httpecho/handlers"
	"github.com/madcarpet/metrics/internal/service/metrics"
	"github.com/stretchr/testify/assert"
)

func TestAsymDecription(t *testing.T) {
	data := `{"id": "test2", "type": "gauge", "value": 15.212}`

	var bodyBuffer bytes.Buffer
	bodyBuffer.Write([]byte(data))

	db := memstorage.NewMemStorage()
	e := echo.New()
	updateSvc := metrics.NewUpdateMetricSvc(db)
	getSvc := metrics.NewGetMetricSvc(db)
	updateHandler := handlers.NewUpdateHandler(updateSvc, getSvc)
	e.POST("/update/", updateHandler.Handle, AsymetricDecrypt("../../../asymetric/test_keys/privkey.pem"))

	req := httptest.NewRequest(http.MethodPost, "http://localhost:8080/update/", &bodyBuffer)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	resp := rec.Result()
	resp.Body.Close()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}
