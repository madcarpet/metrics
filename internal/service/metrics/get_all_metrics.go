package metrics

import (
	"context"

	"github.com/madcarpet/metrics/internal/adapter/storage"
	"github.com/madcarpet/metrics/internal/entity"
)

type GetAllMetricsSvc struct {
	repo storage.Repository
}

func (s *GetAllMetricsSvc) GetAllMetrics(ctx context.Context) ([]entity.Metric, error) {
	return s.repo.GetAllMetrics(ctx)
}

func NewGetAllMetricsSvc(r storage.Repository) *GetAllMetricsSvc {
	return &GetAllMetricsSvc{
		repo: r,
	}
}
