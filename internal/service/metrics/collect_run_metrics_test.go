package metrics

import (
	"context"
	"testing"

	"github.com/madcarpet/metrics/internal/adapter/storage/memstorage"
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

	allocMetric, err := db.GetByNameAndType(context.Background(), "Alloc", entity.Gauge)
	assert.Nil(t, err)
	assert.Equal(t, "Alloc", allocMetric.Name)

	freesMetric, err := db.GetByNameAndType(context.Background(), "Frees", entity.Gauge)
	assert.Nil(t, err)
	assert.Equal(t, "Frees", freesMetric.Name)

	heapMetric, err := db.GetByNameAndType(context.Background(), "HeapAlloc", entity.Gauge)
	assert.Nil(t, err)
	assert.Equal(t, "HeapAlloc", heapMetric.Name)

	gcsysMetric, err := db.GetByNameAndType(context.Background(), "GCSys", entity.Gauge)
	assert.Nil(t, err)
	assert.Equal(t, "GCSys", gcsysMetric.Name)

	randomMetric, err := db.GetByNameAndType(context.Background(), "RandomValue", entity.Gauge)
	assert.Nil(t, err)
	assert.Equal(t, "RandomValue", randomMetric.Name)

	pollCountMetric, err := db.GetByNameAndType(context.Background(), "PollCount", entity.Counter)
	assert.Nil(t, err)
	assert.Equal(t, "PollCount", pollCountMetric.Name)
	assert.Equal(t, 1, int(pollCountMetric.Value))

	collectRunMetricsSvc.Collect(testMetrics)
	pollCountMetric, err = db.GetByNameAndType(context.Background(), "PollCount", entity.Counter)
	assert.Nil(t, err)
	assert.Equal(t, "PollCount", pollCountMetric.Name)
	assert.Equal(t, 2, int(pollCountMetric.Value))
}
