package storage

// type for storage in memory
type MemStorage struct {
	metrics map[string]float64
}

type MetricType int

// func for update key value pairs
func (s *MemStorage) Update(k string, v float64, t MetricType) {
	switch t {
	case 1:
		s.metrics[k] = v
	case 2:
		s.metrics[k] += v
	}

}

// Constructor for MemStorage
func NewMemStorage() *MemStorage {
	return &MemStorage{metrics: make(map[string]float64)}
}
