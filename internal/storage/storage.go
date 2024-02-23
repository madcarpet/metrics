package storage

import (
	"fmt"
	"sync"
)

// interface to user memstorage type
type Repositories interface {
	Update(k string, v float64, t MetricType)
	GetMetric(k string, t MetricType) (string, error)
	GetMetricsAll() (map[string]string, error)
}

// type for storage in memory
type MemStorage struct {
	Gauges   map[string]float64
	Counters map[string]int64
}

// type for identification metric type
type MetricType int

// function to update key value pairs
func (s *MemStorage) Update(k string, v float64, t MetricType) {
	var mutex sync.Mutex
	mutex.Lock()
	switch t {
	case 1:
		s.Gauges[k] = v
	case 2:
		s.Counters[k] += int64(v)
	}
	mutex.Unlock()

}

// function to get metric by key and type
func (s MemStorage) GetMetric(k string, t MetricType) (string, error) {
	switch t {
	case 1:
		metric, ok := s.Gauges[k]
		if !ok {
			return "", fmt.Errorf("ERROR: %v", "No such metric")
		}
		return fmt.Sprint(metric), nil
	default:
		metric, ok := s.Counters[k]
		if !ok {
			return "", fmt.Errorf("ERROR: %v", "No such metric")
		}
		return fmt.Sprint(metric), nil
	}
}

// function to get all metrics
func (s MemStorage) GetMetricsAll() (map[string]string, error) {
	allmetrics := make(map[string]string)
	if len(s.Gauges) > 0 {
		for k, v := range s.Gauges {
			allmetrics[k] = fmt.Sprint(v)
		}
	}
	if len(s.Counters) > 0 {
		for k, v := range s.Counters {
			allmetrics[k] = fmt.Sprint(v)
		}
	}
	return allmetrics, nil
}

// constructor for MemStorage
func NewMemStorage() *MemStorage {
	return &MemStorage{Gauges: make(map[string]float64), Counters: make(map[string]int64)}
}
