package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/madcarpet/metrics/internal/entity"
	"github.com/madcarpet/metrics/internal/models"
)

type repository interface {
	GetAllMetrics() []entity.Metric
}

type reporter struct {
	serverAddress string
	repo          repository
}

func NewReporter(sa string, r repository) *reporter {
	return &reporter{serverAddress: sa, repo: r}
}

func (r *reporter) ReportMetrics() error {
	var reqData models.Metrics
	var body bytes.Buffer
	metrics := r.repo.GetAllMetrics()
	for _, metric := range metrics {
		switch metric.Type {
		case entity.Gauge:
			reqData = models.Metrics{
				ID:    metric.Name,
				MType: "gauge",
				Value: &metric.Value,
			}
		case entity.Counter:
			metricDelta := int64(metric.Value)
			reqData = models.Metrics{
				ID:    metric.Name,
				MType: "counter",
				Delta: &metricDelta,
			}
		}
		jsonBody, err := json.Marshal(&reqData)
		if err != nil {
			return fmt.Errorf("report encoding error: %s", err)
		}
		_, err = body.Write(jsonBody)
		if err != nil {
			return fmt.Errorf("json body writing to buffer error: %s", err)
		}
		url := fmt.Sprintf("http://%v/update/", r.serverAddress)
		resp, err := http.Post(url, "application/json", &body)
		if err != nil {
			return fmt.Errorf("http error: %s", err)
		}
		defer resp.Body.Close()
	}
	return nil
}
