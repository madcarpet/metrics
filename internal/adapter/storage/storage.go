package storage

import (
	"context"

	"github.com/madcarpet/metrics/internal/entity"
)

type Repository interface {
	GetByNameAndType(ctx context.Context, n string, t int64) (entity.Metric, error)
	UpdateMetric(ctx context.Context, m entity.Metric) error
	GetAllMetrics(ctx context.Context) []entity.Metric
	ExportToFile() error
	ImportFromFile() error
	Close() error
	IsConnected(ctx context.Context) error
}
