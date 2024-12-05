package storage

import (
	"context"

	"github.com/madcarpet/metrics/internal/entity"
)

// Repository interface with methods for managing storage.
type Repository interface {
	GetByNameAndType(ctx context.Context, n string, t int64) (entity.Metric, error)
	UpdateMetric(ctx context.Context, m entity.Metric) error
	GetAllMetrics(ctx context.Context) ([]entity.Metric, error)
	ExportToFile(ctx context.Context) error
	ImportFromFile(ctx context.Context) error
	Close() error
	IsConnected(ctx context.Context) error
	UpdateMetrics(ctx context.Context, mcs []entity.Metric) error
}
