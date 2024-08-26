package metrics

import (
	"github.com/madcarpet/metrics/internal/adapter/storage"
	"github.com/madcarpet/metrics/internal/entity"
)

type GetMetricSvc struct {
	repo storage.Repository
}

func (s *GetMetricSvc) GetMetric(n string, t int64) (entity.Metric, error) {
	metric, err := s.repo.GetByNameAndType(n, t)
	if err != nil {
		return entity.Metric{}, err
	}
	return metric, nil
}

func NewGetMetricSvc(r storage.Repository) *GetMetricSvc {
	return &GetMetricSvc{
		repo: r,
	}
}
