package entity

const (
	Gauge int64 = 1 + iota
	Counter
)

type Metric struct {
	Type  int64   `json:"type"`
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}
