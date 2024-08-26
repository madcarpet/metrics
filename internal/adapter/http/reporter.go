package http

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/madcarpet/metrics/internal/adapter/storage"
	"github.com/madcarpet/metrics/internal/entity"
	"github.com/madcarpet/metrics/internal/models"
)

type reporter struct {
	serverAddress string
	repo          storage.Repository
}

func NewReporter(sa string, r storage.Repository) *reporter {
	return &reporter{serverAddress: sa, repo: r}
}

func (r *reporter) ReportMetrics() error {
	var finalReqData []models.Metrics
	var body bytes.Buffer
	metrics := r.repo.GetAllMetrics(context.Background())
	for _, metric := range metrics {
		var reqData models.Metrics
		switch metric.Type {
		case entity.Gauge:
			mName := metric.Name
			mVal := metric.Value
			reqData = models.Metrics{
				ID:    mName,
				MType: "gauge",
				Value: &mVal,
			}
		case entity.Counter:
			mName := metric.Name
			metricDelta := int64(metric.Value)
			reqData = models.Metrics{
				ID:    mName,
				MType: "counter",
				Delta: &metricDelta,
			}
		}
		finalReqData = append(finalReqData, reqData)
	}
	jsonBody, err := json.Marshal(&finalReqData)
	if err != nil {
		return fmt.Errorf("report encoding error: %s", err)
	}

	gzBodyWriter := gzip.NewWriter(&body)
	_, err = gzBodyWriter.Write(jsonBody)
	if err != nil {
		return fmt.Errorf("json body compression to buffer error: %s", err)
	}
	if err := gzBodyWriter.Close(); err != nil {
		return fmt.Errorf("gzip writer closing error: %s", err)
	}
	url := fmt.Sprintf("http://%v/updates/", r.serverAddress)
	req, err := http.NewRequest(http.MethodPost, url, &body)
	if err != nil {
		return fmt.Errorf("request formation error: %s", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Accept-Encoding", "gzip")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("http request sendig error: %s", err)
	}
	defer resp.Body.Close()
	return nil
}
