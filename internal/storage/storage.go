package storage

import "fmt"

// type for storage in memory
type MemStorage struct {
	Gauges   map[string]float64
	Counters map[string]int64
}

type MetricType int

// func for update key value pairs
func (s *MemStorage) Update(k string, v float64, t MetricType) {
	switch t {
	case 1:
		s.Gauges[k] = v
	case 2:
		s.Counters[k] += int64(v)
	}

}

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

// Constructor for MemStorage
func NewMemStorage() *MemStorage {
	return &MemStorage{Gauges: make(map[string]float64), Counters: make(map[string]int64)}
}
