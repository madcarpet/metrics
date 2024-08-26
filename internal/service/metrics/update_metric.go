package metrics

import (
	"fmt"

	"github.com/madcarpet/metrics/internal/entity"
)

type updateMetricRepo interface {
	UpdateMetric(m entity.Metric) error
}

type UpdateMetricSvc struct {
	repo updateMetricRepo
}

func (s *UpdateMetricSvc) UpdateMetric(m entity.Metric) error {
	err := s.repo.UpdateMetric(m)
	if err != nil {
		return fmt.Errorf("error while updating metric: %v", err)
	}
	return nil
}

func NewUpdateMetricSvc(r updateMetricRepo) *UpdateMetricSvc {
	return &UpdateMetricSvc{
		repo: r,
	}
}
