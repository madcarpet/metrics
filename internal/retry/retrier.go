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
	interval time.Duration
	fn       func(ctx context.Context) error
}

func (r retrier) Retry(ctx context.Context) error {
	err := r.fn(ctx)
	if err == nil {
		return nil
	}
	var connErr net.Error
	if errors.As(err, &connErr) {
		for i := 1; i < r.retries; i++ {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				logger.Log.Info("Retry")
				logger.Log.Info(connErr.Error())
				err = r.fn(ctx)
				if err == nil {
					return nil
				}
				if i != r.retries {
					time.Sleep(r.interval)
				}
			}
		}
	}
	return err
}

func NewRetrier(r int, i int, fn func(ctx context.Context) error) *retrier {
	interval := time.Duration(i) * time.Second
	return &retrier{r, interval, fn}
}
