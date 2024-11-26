package dispenser

import (
	"context"
	"time"

	"github.com/madcarpet/metrics/internal/adapter/storage"
	"github.com/madcarpet/metrics/internal/entity"
)

// metricsDispenser - structure for metric dispenser.
// Keep strorage and out channel.
type metricsDispenser struct {
	repo  storage.Repository
	outCh chan<- []entity.Metric
}

// NewMetricDispenser - constructor for metricsDispenser.
func NewMetricDispenser(r storage.Repository, ch chan []entity.Metric) *metricsDispenser {
	return &metricsDispenser{
		repo:  r,
		outCh: ch,
	}
}

// Dispense - function to collect metrics and dispense them to channel.
func (d *metricsDispenser) Dispense(ctx context.Context, ri int64) error {
	tick := time.NewTicker(time.Duration(ri) * time.Second)
	for {
		metrics, err := d.repo.GetAllMetrics(ctx)
		if err != nil {
			return err
		}
		d.outCh <- metrics
		<-tick.C
	}
}
