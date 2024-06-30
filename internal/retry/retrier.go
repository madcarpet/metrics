package retry

import (
	"context"
	"errors"
	"net"
	"time"

	"github.com/madcarpet/metrics/internal/logger"
)

const (
	DefaultRetry = 3
	Interval2s   = 2
)

type retrier struct {
	retries  int
	interval []time.Duration
	fn       func(ctx context.Context) error
}

func (r retrier) Retry(ctx context.Context) error {
	var err error
	for i := 0; i < r.retries; i++ {
		err = r.fn(ctx)
		if err == nil {
			return nil
		}
		var connErr net.Error
		if errors.As(err, &connErr) {
			logger.Log.Info("Retry")
			logger.Log.Info(connErr.Error())
			if i < len(r.interval) {
				time.Sleep(r.interval[i])
			}
		}
	}
	return err
}

func NewRetrier(r int, i int, fn func(ctx context.Context) error) *retrier {
	var intervals []time.Duration
	switch i {
	default:
		intervals = []time.Duration{1 * time.Second, 3 * time.Second, 5 * time.Second}
	}
	return &retrier{r, intervals, fn}
}
