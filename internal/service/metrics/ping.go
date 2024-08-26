package metrics

import (
	"context"

	"github.com/madcarpet/metrics/internal/adapter/storage"
)

type PingSvc struct {
	repo storage.Repository
}

func (s *PingSvc) Ping(ctx context.Context) error {
	return s.repo.IsConnected(ctx)

}

func NewPingSvc(r storage.Repository) *PingSvc {
	return &PingSvc{repo: r}
}
