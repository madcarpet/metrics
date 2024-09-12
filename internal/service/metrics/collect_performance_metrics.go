package metrics

import (
	"context"
	"fmt"

	"github.com/madcarpet/metrics/internal/adapter/storage"
	"github.com/madcarpet/metrics/internal/entity"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
)

type PerfMetricCollectorSvc struct {
	repo storage.Repository
}

func NewPerCollectorSvc(r storage.Repository) *PerfMetricCollectorSvc {
	return &PerfMetricCollectorSvc{repo: r}
}

func (pc *PerfMetricCollectorSvc) Collect(mn []string) error {
	var collecdetMetrics []entity.Metric
	m, err := mem.VirtualMemory()
	if err != nil {
		return err
	}
	mTotalMetric := entity.Metric{
		Type:  entity.Gauge,
		Name:  mn[0],
		Value: float64(m.Total),
	}
	mFreeMetric := entity.Metric{
		Type:  entity.Gauge,
		Name:  mn[1],
		Value: float64(m.Free),
	}
	collecdetMetrics = append(collecdetMetrics, mTotalMetric, mFreeMetric)
	cpus, err := cpu.Percent(5, true)
	if err != nil {
		return err
	}
	for i, c := range cpus {
		cpuMetric := entity.Metric{
			Type:  entity.Gauge,
			Name:  fmt.Sprintf("%s%d", mn[2], i+1),
			Value: c,
		}
		collecdetMetrics = append(collecdetMetrics, cpuMetric)
	}
	err = pc.repo.UpdateMetrics(context.TODO(), collecdetMetrics)
	if err != nil {
		return err
	}
	return nil
}
