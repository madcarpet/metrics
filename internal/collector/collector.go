package collector

import (
	"math/rand"
	"runtime"
)

type Metrics struct {
	Gauge   map[string]float64
	Counter map[string]int64
}

func NewMetrics() Metrics {
	return Metrics{Gauge: make(map[string]float64), Counter: make(map[string]int64)}

}

func (m *Metrics) Collect() {
	st := runtime.MemStats{}
	runtime.ReadMemStats(&st)
	m.Gauge["Alloc"] = float64(st.Alloc)
	m.Gauge["BuckHashSys"] = float64(st.BuckHashSys)
	m.Gauge["Frees"] = float64(st.Frees)
	m.Gauge["GCCPUFraction"] = float64(st.GCCPUFraction)
	m.Gauge["GCSys"] = float64(st.GCSys)
	m.Gauge["HeapAlloc"] = float64(st.HeapAlloc)
	m.Gauge["HeapIdle"] = float64(st.HeapIdle)
	m.Gauge["HeapInuse"] = float64(st.HeapInuse)
	m.Gauge["HeapObjects"] = float64(st.HeapObjects)
	m.Gauge["HeapReleased"] = float64(st.HeapReleased)
	m.Gauge["HeapSys"] = float64(st.HeapSys)
	m.Gauge["LastGC"] = float64(st.LastGC)
	m.Gauge["Lookups"] = float64(st.Lookups)
	m.Gauge["MCacheInuse"] = float64(st.MCacheInuse)
	m.Gauge["MCacheSys"] = float64(st.MCacheSys)
	m.Gauge["MSpanInuse"] = float64(st.MSpanInuse)
	m.Gauge["MSpanSys"] = float64(st.MSpanSys)
	m.Gauge["Mallocs"] = float64(st.Mallocs)
	m.Gauge["NextGC"] = float64(st.NextGC)
	m.Gauge["NumForcedGC"] = float64(st.NumForcedGC)
	m.Gauge["NumGC"] = float64(st.NumGC)
	m.Gauge["OtherSys"] = float64(st.OtherSys)
	m.Gauge["PauseTotalNs"] = float64(st.PauseTotalNs)
	m.Gauge["StackInuse"] = float64(st.StackInuse)
	m.Gauge["StackSys"] = float64(st.StackSys)
	m.Gauge["Sys"] = float64(st.Sys)
	m.Gauge["TotalAlloc"] = float64(st.TotalAlloc)
	m.Gauge["RandomValue"] = rand.Float64()
	m.Counter["PollCount"] += 1
}

func (m Metrics) GetGauge() map[string]float64 {
	return m.Gauge
}

func (m Metrics) GetCounter() map[string]int64 {
	return m.Counter
}
