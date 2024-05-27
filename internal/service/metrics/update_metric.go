package metrics

import (
	"context"
	"fmt"

	"github.com/madcarpet/metrics/internal/adapter/storage"
	"github.com/madcarpet/metrics/internal/entity"
)

type UpdateMetricSvc struct {
	repo storage.Repository
}

func (s *UpdateMetricSvc) UpdateMetric(ctx context.Context, m entity.Metric) error {
	err := s.repo.UpdateMetric(ctx, m)
	if err != nil {
		return fmt.Errorf("error while updating metric: %v", err)
	}
	return nil
}

func NewUpdateMetricSvc(r storage.Repository) *UpdateMetricSvc {
	return &UpdateMetricSvc{
		repo: r,
	}
}
