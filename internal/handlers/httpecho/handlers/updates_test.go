package handlers

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/madcarpet/metrics/internal/entity"
	"github.com/madcarpet/metrics/internal/mocks"
	"github.com/madcarpet/metrics/internal/service/metrics"
	"github.com/stretchr/testify/assert"
)

func TestUpdatesHandler(t *testing.T) {

	type want struct {
		code        int
		contentType string
		respBody    string
	}
	tests := []struct {
		name           string
		method         string
		reqContentType string
		reqBody        string
		url            string
		toStorageData  []entity.Metric
		want           want
	}{
		{
			name:           "Test add new set of parameters",
			url:            "http://localhost:8080/updates/",
			method:         http.MethodPost,
			reqContentType: "application/json",
			reqBody:        `[{"id":"testUpdatesGA1","type":"gauge","value":500},{"id":"testUpdatesGA1","type":"gauge","value":1000},{"id":"testUpdatesGA1","type":"gauge","value":3000},{"id":"testUpdatesCO1","type":"counter","delta":10},{"id":"testUpdatesCO1","type":"counter","delta":15},{"id":"testUpdatesCO1","type":"counter","delta":25}]`,
			toStorageData: []entity.Metric{
				{
					Name:  "testUpdatesGA1",
					Type:  1,
					Value: 3000,
				},
				{
					Name:  "testUpdatesCO1",
					Type:  2,
					Value: 50,
				},
			},
			want: want{
				code:        http.StatusOK,
				contentType: "application/json",
				respBody:    `[{"id":"testUpdatesGA1","type":"gauge","value":3000},{"id":"testUpdatesCO1","type":"counter","delta":50}]`,
			},
		},
		{
			name:           "Test empty ID",
			url:            "http://localhost:8080/updates/",
			method:         http.MethodPost,
			reqContentType: "application/json",
			reqBody:        `[{"id":"","type":"gauge","value":500},{"id":"testUpdatesGA1","type":"gauge","value":1000},{"id":"testUpdatesGA1","type":"gauge","value":3000},{"id":"testUpdatesCO1","type":"counter","delta":10},{"id":"testUpdatesCO1","type":"counter","delta":15},{"id":"testUpdatesCO1","type":"counter","delta":25}]`,
			toStorageData: []entity.Metric{
				{
					Name:  "testUpdatesGA1",
					Type:  1,
					Value: 3000,
				},
				{
					Name:  "testUpdatesCO1",
					Type:  2,
					Value: 50,
				},
			},
			want: want{
				code:        http.StatusNotFound,
				contentType: "text/plain; charset=UTF-8",
				respBody:    "Metric name not found for some metrics",
			},
		},
		{
			name:           "Test no ID",
			url:            "http://localhost:8080/updates/",
			method:         http.MethodPost,
			reqContentType: "application/json",
			reqBody:        `[{"type":"gauge","value":500},{"id":"testUpdatesGA1","type":"gauge","value":1000},{"id":"testUpdatesGA1","type":"gauge","value":3000},{"id":"testUpdatesCO1","type":"counter","delta":10},{"id":"testUpdatesCO1","type":"counter","delta":15},{"id":"testUpdatesCO1","type":"counter","delta":25}]`,
			toStorageData: []entity.Metric{
				{
					Name:  "testUpdatesGA1",
					Type:  1,
					Value: 3000,
				},
				{
					Name:  "testUpdatesCO1",
					Type:  2,
					Value: 50,
				},
			},
			want: want{
				code:        http.StatusNotFound,
				contentType: "text/plain; charset=UTF-8",
				respBody:    "Metric name not found for some metrics",
			},
		},
		{
			name:           "Bad gauge value",
			url:            "http://localhost:8080/updates/",
			method:         http.MethodPost,
			reqContentType: "application/json",
			reqBody:        `[{"id":"testUpdatesGA1","type":"gauge","value":"500d"},{"id":"testUpdatesGA1","type":"gauge","value":1000},{"id":"testUpdatesGA1","type":"gauge","value":3000},{"id":"testUpdatesCO1","type":"counter","delta":10},{"id":"testUpdatesCO1","type":"counter","delta":15},{"id":"testUpdatesCO1","type":"counter","delta":25}]`,
			toStorageData: []entity.Metric{
				{
					Name:  "testUpdatesGA1",
					Type:  1,
					Value: 3000,
				},
				{
					Name:  "testUpdatesCO1",
					Type:  2,
					Value: 50,
				},
			},
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain; charset=UTF-8",
				respBody:    "Bad request",
			},
		},
		{
			name:           "Guge and delta mismatch",
			url:            "http://localhost:8080/updates/",
			method:         http.MethodPost,
			reqContentType: "application/json",
			reqBody:        `[{"id":"testUpdatesGA1","type":"gauge","value":"500d"},{"id":"testUpdatesGA1","type":"gauge","value":1000},{"id":"testUpdatesGA1","type":"gauge","value":3000},{"id":"testUpdatesCO1","type":"counter","delta":10},{"id":"testUpdatesCO1","type":"counter","delta":15},{"id":"testUpdatesCO1","type":"counter","delta":25}]`,
			toStorageData: []entity.Metric{
				{
					Name:  "testUpdatesGA1",
					Type:  1,
					Value: 3000,
				},
				{
					Name:  "testUpdatesCO1",
					Type:  2,
					Value: 50,
				},
			},
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain; charset=UTF-8",
				respBody:    "Bad request",
			},
		},
		{
			name:           "Bad delta value",
			url:            "http://localhost:8080/updates/",
			method:         http.MethodPost,
			reqContentType: "application/json",
			reqBody:        `[{"id":"testUpdatesGA1","type":"gauge","value":"500"},{"id":"testUpdatesGA1","type":"gauge","value":1000},{"id":"testUpdatesGA1","type":"gauge","value":3000},{"id":"testUpdatesCO1","type":"counter","delta":"1s"},{"id":"testUpdatesCO1","type":"counter","delta":15},{"id":"testUpdatesCO1","type":"counter","delta":25}]`,
			toStorageData: []entity.Metric{
				{
					Name:  "testUpdatesGA1",
					Type:  1,
					Value: 3000,
				},
				{
					Name:  "testUpdatesCO1",
					Type:  2,
					Value: 50,
				},
			},
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain; charset=UTF-8",
				respBody:    "Bad request",
			},
		},
		{
			name:           "Missed value",
			url:            "http://localhost:8080/updates/",
			method:         http.MethodPost,
			reqContentType: "application/json",
			reqBody:        `[{"id":"testUpdatesGA1","type":"gauge","value":"500"},{"id":"testUpdatesGA1","type":"gauge"},{"id":"testUpdatesGA1","type":"gauge","value":3000},{"id":"testUpdatesCO1","type":"counter","delta":10},{"id":"testUpdatesCO1","type":"counter","delta":15},{"id":"testUpdatesCO1","type":"counter","delta":25}]`,
			toStorageData: []entity.Metric{
				{
					Name:  "testUpdatesGA1",
					Type:  1,
					Value: 3000,
				},
				{
					Name:  "testUpdatesCO1",
					Type:  2,
					Value: 50,
				},
			},
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain; charset=UTF-8",
				respBody:    "Bad request",
			},
		},
		{
			name:           "Missed delta",
			url:            "http://localhost:8080/updates/",
			method:         http.MethodPost,
			reqContentType: "application/json",
			reqBody:        `[{"id":"testUpdatesGA1","type":"gauge","value":"500d"},{"id":"testUpdatesGA1","type":"gauge","value":1000},{"id":"testUpdatesGA1","type":"gauge","value":3000},{"id":"testUpdatesCO1","type":"counter","delta":10},{"id":"testUpdatesCO1","type":"counter","delta":15},{"id":"testUpdatesCO1","type":"counter"}]`,
			toStorageData: []entity.Metric{
				{
					Name:  "testUpdatesGA1",
					Type:  1,
					Value: 3000,
				},
				{
					Name:  "testUpdatesCO1",
					Type:  2,
					Value: 50,
				},
			},
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain; charset=UTF-8",
				respBody:    "Bad request",
			},
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	s := mocks.NewMockRepository(ctrl)
	s.EXPECT().UpdateMetrics(context.Background(), tests[0].toStorageData).Return(nil)

	e := echo.New()
	updatesSvc := metrics.NewUpdateMetricsSvc(s)
	getSvc := metrics.NewGetMetricSvc(s)
	updatesHandler := NewUpdatesHandler(updatesSvc, getSvc)
	e.POST("/updates/", updatesHandler.Handle)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fmt.Println(test.toStorageData)
			var bodyBuffer bytes.Buffer
			bodyBuffer.Write([]byte(test.reqBody))
			req := httptest.NewRequest(test.method, test.url, &bodyBuffer)
			req.Header.Set("Content-Type", test.reqContentType)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)
			defer rec.Result().Body.Close()
			assert.Equal(t, test.want.code, rec.Code)
			assert.Equal(t, test.want.contentType, rec.Header().Get("Content-Type"))
			resRespBody := strings.TrimRight(rec.Body.String(), "\n")
			assert.Equal(t, test.want.respBody, resRespBody)
		},
		)

	}
}
