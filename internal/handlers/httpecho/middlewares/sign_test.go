package middlewares

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/madcarpet/metrics/internal/adapter/storage/memstorage"
	"github.com/madcarpet/metrics/internal/handlers/httpecho/handlers"
	"github.com/madcarpet/metrics/internal/service/metrics"
	"github.com/stretchr/testify/assert"
)

func TestSignMiddleware(t *testing.T) {
	data := `{"id": "test2", "type": "gauge", "value": 15.212}`
	respData := "{\"id\":\"test2\",\"type\":\"gauge\",\"value\":15.212}\n"
	reqSign := signatory([]byte(data), "topsecret")

	os.Setenv("SECRET_KEY", "topsecret")
	defer os.Clearenv()

	var bodyBuffer bytes.Buffer
	var respBuffer bytes.Buffer
	bodyBuffer.Write([]byte(data))
	respBuffer.Write([]byte(respData))

	fmt.Println(respBuffer.Bytes())
	respSign := signatory(respBuffer.Bytes(), "topsecret")

	db := memstorage.NewMemStorage()
	e := echo.New()
	updateSvc := metrics.NewUpdateMetricSvc(db)
	getSvc := metrics.NewGetMetricSvc(db)
	updateHandler := handlers.NewUpdateHandler(updateSvc, getSvc)
	e.POST("/update/", updateHandler.Handle, SignData)

	req := httptest.NewRequest(http.MethodPost, "http://localhost:8080/update/", &bodyBuffer)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("HashSHA256", reqSign)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	resp := rec.Result()
	resp.Body.Close()

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, respSign, rec.Header().Get("HashSHA256"))
	respBody, _ := io.ReadAll(resp.Body)
	assert.Equal(t, respData, string(respBody))

}

func TestSignatory(t *testing.T) {
	hash := signatory([]byte("hello"), "hello")
	assert.Equal(t, "b270147ff516860aafb4f52818cb149defde2fac57b451fa3b053d24e88b915a", hash)
}
