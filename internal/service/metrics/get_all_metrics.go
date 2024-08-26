package metrics

import "github.com/madcarpet/metrics/internal/entity"

type getAllMetricsRepo interface {
	GetAllMetrics() []entity.Metric
}

type GetAllMetricsSvc struct {
	repo getAllMetricsRepo
}

func (s *GetAllMetricsSvc) GetAllMetrics() []entity.Metric {
	return s.repo.GetAllMetrics()
}

func NewGetAllMetricsSvc(r getAllMetricsRepo) *GetAllMetricsSvc {
	return &GetAllMetricsSvc{
		repo: r,
	}
}
