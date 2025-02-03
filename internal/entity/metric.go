// Package entity - contains app entity structures.
package entity

// Constans for differentiating between gauge and counter metric types.
const (
	Gauge int64 = 1 + iota
	Counter
)

// Metric - main structure for metric.
type Metric struct {
	Type  int64   `json:"type"`
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}
