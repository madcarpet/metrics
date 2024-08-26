package storage

import (
	"context"

	"github.com/madcarpet/metrics/internal/entity"
)

type Repository interface {
	GetByNameAndType(n string, t int64) (entity.Metric, error)
	UpdateMetric(m entity.Metric) error
	GetAllMetrics() []entity.Metric
	ExportToFile() error
	ImportFromFile() error
	Close() error
	IsConnected(ctx context.Context) error
}
