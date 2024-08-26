package metrics

import (
	"github.com/madcarpet/metrics/internal/adapter/storage"
	"github.com/madcarpet/metrics/internal/entity"
)

type GetAllMetricsSvc struct {
	repo storage.Repository
}

func (s *GetAllMetricsSvc) GetAllMetrics() []entity.Metric {
	return s.repo.GetAllMetrics()
}

func NewGetAllMetricsSvc(r storage.Repository) *GetAllMetricsSvc {
	return &GetAllMetricsSvc{
		repo: r,
	}
}
