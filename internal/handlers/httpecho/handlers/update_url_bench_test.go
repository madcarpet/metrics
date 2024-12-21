package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/madcarpet/metrics/internal/adapter/storage/memstorage"
	"github.com/madcarpet/metrics/internal/service/metrics"
)

func BenchmarkUpdateUrlHandler(b *testing.B) {
	db := memstorage.NewMemStorage()
	e := echo.New()
	updateSvc := metrics.NewUpdateMetricSvc(db)
	updateURLHandler := NewUpdateURLHandler(updateSvc)
	e.POST("/update/:type/:name/:value", updateURLHandler.Handle)

	reqG := httptest.NewRequest(http.MethodPost, "http://localhost:8080/update/gauge/TestGauge1/1.21511", nil)
	reqC := httptest.NewRequest(http.MethodPost, "http://localhost:8080/update/counter/TestCounter1/188", nil)
	recG := httptest.NewRecorder()
	recC := httptest.NewRecorder()
	defer recG.Result().Body.Close()
	defer recC.Result().Body.Close()
	b.ResetTimer()
	b.Run("update url gauge benchmark", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			e.ServeHTTP(recG, reqG)
		}
	})
	b.Run("update url counter benchmark", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			e.ServeHTTP(recC, reqC)
		}
	})

}
