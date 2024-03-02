package metrics

import (
	"testing"

	"github.com/madcarpet/metrics/internal/adapter/memstorage"
	"github.com/madcarpet/metrics/internal/entity"
	"github.com/stretchr/testify/assert"
)

func TestCollectRunMetrics(t *testing.T) {
	db := memstorage.NewMemStorage()
	wrongMetric := []string{"WrongMetricName"}
	testMetrics := []string{"Alloc", "Frees", "HeapAlloc", "GCSys"}
	collectRunMetricsSvc := NewCollectorSvc(db)

	err := collectRunMetricsSvc.Collect(wrongMetric)
	assert.NotNil(t, err)

	err = collectRunMetricsSvc.Collect(testMetrics)
	assert.Nil(t, err)

	allocMetric, err := db.GetByNameAndType("Alloc", entity.Gauge)
	assert.Nil(t, err)
	assert.Equal(t, allocMetric.Name, "Alloc")

	freesMetric, err := db.GetByNameAndType("Frees", entity.Gauge)
	assert.Nil(t, err)
	assert.Equal(t, freesMetric.Name, "Frees")

	heapMetric, err := db.GetByNameAndType("HeapAlloc", entity.Gauge)
	assert.Nil(t, err)
	assert.Equal(t, heapMetric.Name, "HeapAlloc")

	gcsysMetric, err := db.GetByNameAndType("GCSys", entity.Gauge)
	assert.Nil(t, err)
	assert.Equal(t, gcsysMetric.Name, "GCSys")

	randomMetric, err := db.GetByNameAndType("RandomValue", entity.Gauge)
	assert.Nil(t, err)
	assert.Equal(t, randomMetric.Name, "RandomValue")

	pollCountMetric, err := db.GetByNameAndType("PollCount", entity.Counter)
	assert.Nil(t, err)
	assert.Equal(t, pollCountMetric.Name, "PollCount")
	assert.Equal(t, int(pollCountMetric.Value), 1)

	collectRunMetricsSvc.Collect(testMetrics)
	pollCountMetric, err = db.GetByNameAndType("PollCount", entity.Counter)
	assert.Nil(t, err)
	assert.Equal(t, pollCountMetric.Name, "PollCount")
	assert.Equal(t, int(pollCountMetric.Value), 2)
}
