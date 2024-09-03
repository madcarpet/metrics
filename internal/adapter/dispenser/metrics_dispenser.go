package dispenser

import (
	"context"
	"time"

	"github.com/madcarpet/metrics/internal/adapter/storage"
	"github.com/madcarpet/metrics/internal/entity"
)

type metricsDispenser struct {
	repo  storage.Repository
	outCh chan<- []entity.Metric
}

func NewMetricDispenser(r storage.Repository, ch chan []entity.Metric) *metricsDispenser {
	return &metricsDispenser{
		repo:  r,
		outCh: ch,
	}
}

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
