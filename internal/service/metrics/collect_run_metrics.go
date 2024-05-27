package metrics

import (
	"context"
	"fmt"
	"math/rand"
	"reflect"
	"runtime"

	"github.com/madcarpet/metrics/internal/adapter/storage"
	"github.com/madcarpet/metrics/internal/entity"
)

func NewCollectorSvc(r storage.Repository) *MetricCollectorSvc {
	return &MetricCollectorSvc{repo: r}
}

type MetricCollectorSvc struct {
	repo storage.Repository
}

func (c *MetricCollectorSvc) Collect(ms []string) error {
	var metric entity.Metric
	stats := runtime.MemStats{}
	runtime.ReadMemStats(&stats)
	value := reflect.ValueOf(stats)
	for _, key := range ms {
		mType := value.FieldByName(key).Kind()
		switch mType {
		case reflect.Uint64:
			metric.Value = float64(value.FieldByName(key).Uint())
		case reflect.Uint32:
			metric.Value = float64(value.FieldByName(key).Uint())
		case reflect.Float64:
			metric.Value = float64(value.FieldByName(key).Float())
		default:
			return fmt.Errorf("wrong metric type or metric doesn't exist")
		}
		metric.Type = entity.Gauge
		metric.Name = key
		err := c.repo.UpdateMetric(context.Background(), metric)
		if err != nil {
			return fmt.Errorf("couldn't store metric")
		}
	}
	metric.Type = entity.Gauge
	metric.Name = "RandomValue"
	metric.Value = rand.Float64()
	err := c.repo.UpdateMetric(context.Background(), metric)
	if err != nil {
		return fmt.Errorf("couldn't store metric")
	}
	metric.Type = entity.Counter
	metric.Name = "PollCount"
	metric.Value = 1
	err = c.repo.UpdateMetric(context.Background(), metric)
	if err != nil {
		return fmt.Errorf("couldn't store metric")
	}

	return nil
}
