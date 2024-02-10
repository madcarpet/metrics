package collector

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCollector(t *testing.T) {
	m := NewMetrics()
	m.Collect()
	assert.NotEqual(t, 1, m.Counter["PollCount"])
	gaugeTestMap := m.GetGauge()
	if _, err := gaugeTestMap["Alloc"]; err != true {
		t.Errorf("Gauge collecting err")
	}
	counterTestMap := m.GetCounter()
	if _, err := counterTestMap["PollCount"]; err != true {
		t.Errorf("Counter collecting err")
	}

}
