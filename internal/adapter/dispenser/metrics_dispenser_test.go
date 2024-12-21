package dispenser

import (
	"context"
	"testing"
	"time"

	"github.com/madcarpet/metrics/internal/adapter/storage/memstorage"
	"github.com/madcarpet/metrics/internal/entity"
	"github.com/stretchr/testify/assert"
)

func TestDispenser(t *testing.T) {
	storage := memstorage.NewMemStorage()
	valueMetrics := []entity.Metric{
		{Name: "TestGauge1", Type: entity.Gauge, Value: 1.12},
		{Name: "TestGauge2", Type: entity.Gauge, Value: 212.921351351141},
		{Name: "TestGauge3", Type: entity.Gauge, Value: 0},
		{Name: "TestCounter1", Type: entity.Counter, Value: 1},
		{Name: "TestCounter2", Type: entity.Counter, Value: 991112111111},
		{Name: "TestCounter3", Type: entity.Counter, Value: 0},
	}

	for _, metric := range valueMetrics {
		storage.UpdateMetric(context.TODO(), metric)
	}

	outCh := make(chan []entity.Metric, 1)

	dispenser := NewMetricDispenser(storage, outCh)

	go func() {
		err := dispenser.Dispense(context.Background(), 1) // interval of 1 second
		assert.NoError(t, err)
	}()

	select {
	case receivedMetrics := <-outCh:
		assert.Equal(t, valueMetrics, receivedMetrics)
	case <-time.After(2 * time.Second):
		t.Fatalf("Metrics were not received in time")
	}

}
