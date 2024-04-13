package metrics

import (
	"fmt"

	"github.com/madcarpet/metrics/internal/adapter/storage"
	"github.com/madcarpet/metrics/internal/entity"
)

type UpdateMetricSvc struct {
	repo storage.Repository
}

func (s *UpdateMetricSvc) UpdateMetric(m entity.Metric) error {
	err := s.repo.UpdateMetric(m)
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
