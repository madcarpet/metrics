package memstorage

import (
	"context"
	"fmt"
	"sync"

	"github.com/madcarpet/metrics/internal/entity"
)

type MemStorage struct {
	metrics []entity.Metric
	mutex   sync.Mutex
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		metrics: []entity.Metric{},
		mutex:   sync.Mutex{},
	}
}

func (s *MemStorage) GetByNameAndType(ctx context.Context, n string, t int64) (entity.Metric, error) {
	for _, m := range s.metrics {
		if m.Type == t && m.Name == n {
			return m, nil
		}
	}
	return entity.Metric{}, fmt.Errorf("metric %s not found", n)
}

func (s *MemStorage) UpdateMetric(ctx context.Context, m entity.Metric) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if len(s.metrics) > 0 {
		for idx, metric := range s.metrics {
			if metric.Name == m.Name && metric.Type == m.Type {
				switch m.Type {
				case entity.Counter:
					s.metrics[idx].Value += m.Value
				case entity.Gauge:
					s.metrics[idx] = m
				}
				return nil
			}
		}
	}
	s.metrics = append(s.metrics, m)
	return nil
}

func (s *MemStorage) GetAllMetrics(ctx context.Context) ([]entity.Metric, error) {
	return s.metrics, nil
}

func (s *MemStorage) ExportToFile() error {
	return nil
}

func (s *MemStorage) ImportFromFile() error {
	return nil
}

func (s *MemStorage) Close() error {
	return nil
}

func (s *MemStorage) IsConnected(ctx context.Context) error {
	return fmt.Errorf("DB type without connection support")
}

func (s *MemStorage) UpdateMetrics(ctx context.Context, mcs []entity.Metric) error {
	for _, m := range mcs {
		s.UpdateMetric(context.TODO(), m)
	}
	return nil
}
