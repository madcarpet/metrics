package entity

// Constans for differentiating between gauge and counter metric types.
const (
	Gauge int64 = 1 + iota
	Counter
)

// Main structure for metric.
type Metric struct {
	Type  int64   `json:"type"`
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}
