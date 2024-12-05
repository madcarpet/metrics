package handlers

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/madcarpet/metrics/internal/entity"
	"github.com/madcarpet/metrics/internal/mocks"
	"github.com/madcarpet/metrics/internal/service/metrics"
)

func BenchmarkUpdatesHandler(b *testing.B) {
	ctrl := gomock.NewController(b)
	defer ctrl.Finish()
	s := mocks.NewMockRepository(ctrl)
	toStorageData :=
		[]entity.Metric{
			{
				Name:  "testUpdatesGA1",
				Type:  1,
				Value: 3000,
			},
			{
				Name:  "testUpdatesCO1",
				Type:  2,
				Value: 50,
			}}

	s.EXPECT().UpdateMetrics(context.Background(), toStorageData).Return(nil)
	reqBody := `[{"id":"testUpdatesGA1","type":"gauge","value":500},{"id":"testUpdatesGA1","type":"gauge","value":1000},{"id":"testUpdatesGA1","type":"gauge","value":3000},{"id":"testUpdatesCO1","type":"counter","delta":10},{"id":"testUpdatesCO1","type":"counter","delta":15},{"id":"testUpdatesCO1","type":"counter","delta":25}]`
	e := echo.New()
	updatesSvc := metrics.NewUpdateMetricsSvc(s)
	getSvc := metrics.NewGetMetricSvc(s)
	updatesHandler := NewUpdatesHandler(updatesSvc, getSvc)
	e.POST("/updates/", updatesHandler.Handle)
	var bodyBuffer bytes.Buffer
	bodyBuffer.Write([]byte(reqBody))
	req := httptest.NewRequest(http.MethodPost, "http://localhost:8080/updates/", &bodyBuffer)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	defer rec.Result().Body.Close()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		e.ServeHTTP(rec, req)
	}

}
