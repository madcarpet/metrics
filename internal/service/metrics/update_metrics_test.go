package metrics

import (
	"context"
	"testing"

	"github.com/madcarpet/metrics/internal/adapter/storage/memstorage"
	"github.com/madcarpet/metrics/internal/entity"
	"github.com/stretchr/testify/require"
)

func TestUpdateMetricsSvc(t *testing.T) {
	storage := memstorage.NewMemStorage()
	valueMetrics := []entity.Metric{
		{Name: "TestGauge1", Type: entity.Gauge, Value: 1.12},
		{Name: "TestGauge2", Type: entity.Gauge, Value: 212.921351351141},
		{Name: "TestGauge3", Type: entity.Gauge, Value: 0},
		{Name: "TestCounter1", Type: entity.Counter, Value: 1},
		{Name: "TestCounter2", Type: entity.Counter, Value: 991112111111},
		{Name: "TestCounter3", Type: entity.Counter, Value: 0},
	}
	updSvc := NewUpdateMetricsSvc(storage)
	err := updSvc.UpdateMetrics(context.Background(), valueMetrics)
	require.NoError(t, err)
}
