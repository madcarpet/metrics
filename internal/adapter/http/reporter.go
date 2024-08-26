package http

import (
	"fmt"
	"net/http"

	"github.com/madcarpet/metrics/internal/entity"
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
	var mType string
	metrics := r.repo.GetAllMetrics()
	for _, metric := range metrics {
		switch metric.Type {
		case entity.Gauge:
			mType = "gauge"
		case entity.Counter:
			mType = "counter"
		}
		url := fmt.Sprintf("http://%v/update/%v/%v/%v", r.serverAddress, mType, metric.Name, metric.Value)
		r, err := http.Post(url, "text/plain", nil)
		if err != nil {
			return fmt.Errorf("http error: %s", err)
		}
		defer r.Body.Close()
	}
	return nil
}
