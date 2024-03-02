package metrics

import (
	"testing"

	"github.com/madcarpet/metrics/internal/adapter/memstorage"
	"github.com/madcarpet/metrics/internal/entity"
	"github.com/stretchr/testify/assert"
)

func TestGetMetric(t *testing.T) {
	db := memstorage.NewMemStorage()
	testMetrics := []entity.Metric{
		{Name: "Alloc", Type: entity.Gauge, Value: 2.678488e+06},
		{Name: "BuckHashSys", Type: entity.Gauge, Value: 6160},
		{Name: "GCCPUFraction", Type: entity.Gauge, Value: 0},
		{Name: "RandomValue", Type: entity.Gauge, Value: 0.9233588813342314},
		{Name: "PollCount", Type: entity.Counter, Value: 991112111111},
	}
	for _, metric := range testMetrics {
		db.UpdateMetric(metric)
	}
	getMetricSvc := NewGetMetricSvc(db)
	for _, metric := range testMetrics {
		m, err := getMetricSvc.GetMetric(metric.Name, metric.Type)
		assert.Equal(t, m, metric)
		assert.Nil(t, err)
	}
	_, err := getMetricSvc.GetMetric("fakemetric", entity.Gauge)
	assert.NotNil(t, err)

}
