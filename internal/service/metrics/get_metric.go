package metrics

import "github.com/madcarpet/metrics/internal/entity"

type getMetricRepo interface {
	GetByNameAndType(n string, t int64) (entity.Metric, error)
}

type GetMetricSvc struct {
	repo getMetricRepo
}

func (s *GetMetricSvc) GetMetric(n string, t int64) (entity.Metric, error) {
	metric, err := s.repo.GetByNameAndType(n, t)
	if err != nil {
		return entity.Metric{}, err
	}
	return metric, nil
}

func NewGetMetricSvc(r getMetricRepo) *GetMetricSvc {
	return &GetMetricSvc{
		repo: r,
	}
}
