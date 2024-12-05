package memstorage

import (
	"context"
	"fmt"
	"sync"

	"github.com/madcarpet/metrics/internal/entity"
)

// MemStorage structure for in memory storage.
type MemStorage struct {
	metrics []entity.Metric
	mutex   sync.Mutex
	rmutex  sync.RWMutex
}

// NewMemStorage creates a new MemStorage.
func NewMemStorage() *MemStorage {
	return &MemStorage{
		metrics: []entity.Metric{},
		mutex:   sync.Mutex{},
	}
}

// GetByNameAndType retrive metric by name and type from storage.
func (s *MemStorage) GetByNameAndType(_ context.Context, n string, t int64) (entity.Metric, error) {
	s.rmutex.RLock()
	defer s.rmutex.RUnlock()
	for _, m := range s.metrics {
		if m.Type == t && m.Name == n {
			return m, nil
		}
	}
	return entity.Metric{}, fmt.Errorf("metric %s not found", n)
}

// UpdateMetric updates metric in storage.
func (s *MemStorage) UpdateMetric(_ context.Context, m entity.Metric) error {
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

// GetAllMetrics retrieves all metrics from storage.
func (s *MemStorage) GetAllMetrics(_ context.Context) ([]entity.Metric, error) {
	return s.metrics, nil
}

// ExportToFile not supported, exists just because common interface.
func (s *MemStorage) ExportToFile(_ context.Context) error {
	return nil
}

// ImportFromFile not supported, exists just because common interface.
func (s *MemStorage) ImportFromFile(_ context.Context) error {
	return nil
}

// Close not supported, exists just because common interface.
func (s *MemStorage) Close() error {
	return nil
}

// IsConnected not supported, exists just because common interface.
func (s *MemStorage) IsConnected(_ context.Context) error {
	return fmt.Errorf("storage type without connection support")
}

// UpdateMetrics updates metrics (BUTCH)
func (s *MemStorage) UpdateMetrics(ctx context.Context, mcs []entity.Metric) error {
	for _, m := range mcs {
		s.UpdateMetric(ctx, m)
	}
	return nil
}
