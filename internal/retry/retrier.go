package retry

import (
	"context"
	"errors"
	"fmt"
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
	err := r.fn(ctx)
	if err == nil {
		return nil
	}
	fmt.Println("InRrtry")
	var connErr net.Error
	if errors.As(err, &connErr) {
		for i := 0; i < r.retries; i++ {
			fmt.Println("Retrying")
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				logger.Log.Info("Retry")
				logger.Log.Info(connErr.Error())
				if i < len(r.interval) {
					time.Sleep(r.interval[i])
				}
				err = r.fn(ctx)
				if err == nil {
					return nil
				}
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
