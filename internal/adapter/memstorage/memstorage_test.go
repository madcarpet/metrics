package memstorage

import (
	"testing"

	"github.com/madcarpet/metrics/internal/entity"
	"github.com/stretchr/testify/assert"
)

func TestMemStorage(t *testing.T) {
	storage := NewMemStorage()
	valueMetrics := []entity.Metric{
		{Name: "TestGauge1", Type: entity.Gauge, Value: 1.12},
		{Name: "TestGauge2", Type: entity.Gauge, Value: 212.921351351141},
		{Name: "TestGauge3", Type: entity.Gauge, Value: 0},
		{Name: "TestCounter1", Type: entity.Counter, Value: 1},
		{Name: "TestCounter2", Type: entity.Counter, Value: 991112111111},
		{Name: "TestCounter3", Type: entity.Counter, Value: 0},
	}

	updateMetrics := []entity.Metric{
		{Name: "TestGauge1", Type: entity.Gauge, Value: 1.111111},
		{Name: "TestGauge4", Type: entity.Gauge, Value: 4.44},
		{Name: "TestCounter1", Type: entity.Counter, Value: 99},
		{Name: "TestCounter4", Type: entity.Counter, Value: 4},
	}

	resultMetrics := []entity.Metric{
		{Name: "TestGauge1", Type: entity.Gauge, Value: 1.111111},
		{Name: "TestGauge2", Type: entity.Gauge, Value: 212.921351351141},
		{Name: "TestGauge3", Type: entity.Gauge, Value: 0},
		{Name: "TestCounter1", Type: entity.Counter, Value: 99},
		{Name: "TestCounter2", Type: entity.Counter, Value: 991112111111},
		{Name: "TestCounter3", Type: entity.Counter, Value: 0},
		{Name: "TestGauge4", Type: entity.Gauge, Value: 4.44},
		{Name: "TestCounter4", Type: entity.Counter, Value: 4},
	}

	//Initial dp filling
	for _, metric := range valueMetrics {
		storage.UpdateMetric(metric)
	}
	//Get metrics testing
	for _, metric := range valueMetrics {
		testMetric, _ := storage.GetByNameAndType(metric.Name, metric.Type)
		assert.Equal(t, testMetric.Name, metric.Name)
		assert.Equal(t, testMetric.Type, metric.Type)
		assert.Equal(t, testMetric.Value, metric.Value)
	}
	//Update metrics
	for _, metric := range updateMetrics {
		storage.UpdateMetric(metric)
	}
	//Updated metric testing
	for _, metric := range updateMetrics {
		testMetric, _ := storage.GetByNameAndType(metric.Name, metric.Type)
		assert.Equal(t, testMetric.Name, metric.Name)
		assert.Equal(t, testMetric.Type, metric.Type)
		assert.Equal(t, testMetric.Value, metric.Value)
	}

	//Test not created metric
	_, err := storage.GetByNameAndType("nometric", entity.Counter)
	assert.NotNil(t, err)

	//Test all metrics stored
	assert.Equal(t, len(storage.metrics), 8)

	//Test get all metrics
	assert.Equal(t, storage.metrics, resultMetrics)

}
