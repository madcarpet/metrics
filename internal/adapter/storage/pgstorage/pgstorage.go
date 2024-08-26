package pgstorage

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/madcarpet/metrics/internal/entity"
)

type PGStorage struct {
	DB *sql.DB
}

func NewPGStorage(Params string) (*PGStorage, error) {
	db, err := sql.Open("pgx", Params)
	if err != nil {
		return nil, err
	}
	return &PGStorage{
		DB: db,
	}, nil
}

func (s *PGStorage) IsConnected(ctx context.Context) error {
	subCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	if err := s.DB.PingContext(subCtx); err != nil {
		return err
	}
	return nil
}

func (s *PGStorage) Close() error {
	return s.DB.Close()
}

func (s *PGStorage) GetByNameAndType(n string, t int64) (entity.Metric, error) {
	return entity.Metric{}, nil
}

func (s *PGStorage) UpdateMetric(m entity.Metric) error {
	return nil
}

func (s *PGStorage) GetAllMetrics() []entity.Metric {
	return nil
}

func (s *PGStorage) ExportToFile() error {
	return nil
}

func (s *PGStorage) ImportFromFile() error {
	return nil
}
