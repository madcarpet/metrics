package metrics

import (
	"testing"

	"github.com/madcarpet/metrics/internal/adapter/memstorage"
	"github.com/madcarpet/metrics/internal/entity"
	"github.com/stretchr/testify/assert"
)

func TestUpdateMetric(t *testing.T) {
	db := memstorage.NewMemStorage()
	testMetrics := []entity.Metric{
		{Name: "Alloc", Type: entity.Gauge, Value: 2.678488e+06},
		{Name: "BuckHashSys", Type: entity.Gauge, Value: 6160},
		{Name: "GCCPUFraction", Type: entity.Gauge, Value: 0},
		{Name: "RandomValue", Type: entity.Gauge, Value: 0.9233588813342314},
		{Name: "PollCount", Type: entity.Counter, Value: 991112111111},
	}
	uptateMetricSvc := NewUpdateMetricSvc(db)
	for _, metric := range testMetrics {
		uptateMetricSvc.UpdateMetric(metric)
	}
	assert.Equal(t, db.GetAllMetrics(), testMetrics)
}
