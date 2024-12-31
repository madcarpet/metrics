package metrics

import (
	"context"
	"fmt"

	"github.com/madcarpet/metrics/internal/adapter/storage"
	"github.com/madcarpet/metrics/internal/entity"
)

type UpdateMetricsSvc struct {
	repo storage.Repository
}

func (s *UpdateMetricsSvc) UpdateMetrics(ctx context.Context, mcs []entity.Metric) error {
	err := s.repo.UpdateMetrics(ctx, mcs)
	if err != nil {
		return fmt.Errorf("error while updating metrics: %v", err)
	}
	return nil
}

func NewUpdateMetricsSvc(r storage.Repository) *UpdateMetricsSvc {
	return &UpdateMetricsSvc{
		repo: r,
	}
}
