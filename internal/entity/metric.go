package entity

const (
	Gauge int64 = 1 + iota
	Counter
)

type Metric struct {
	Type  int64
	Name  string
	Value float64
}
