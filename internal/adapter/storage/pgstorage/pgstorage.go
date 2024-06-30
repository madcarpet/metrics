package pgstorage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/madcarpet/metrics/internal/entity"
	"github.com/madcarpet/metrics/internal/retry"
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
	subCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	r := retry.NewRetrier(retry.DefaultRetry, retry.Interval2s, func(ctx context.Context) error {
		if err := s.DB.PingContext(subCtx); err != nil {
			fmt.Println(err)
			return err
		}
		return nil
	})
	return r.Retry(subCtx)
}

func (s *PGStorage) Close() error {
	return s.DB.Close()
}

func (s *PGStorage) GetByNameAndType(ctx context.Context, n string, t int64) (entity.Metric, error) {
	var m entity.Metric
	r := retry.NewRetrier(retry.DefaultRetry, retry.Interval2s, func(ctx context.Context) error {
		tx, err := s.DB.Begin()
		if err != nil {
			return err
		}
		resp := tx.QueryRowContext(ctx, "SELECT type,name,value FROM metrics WHERE type = $1 AND name = $2", t, n)

		err = resp.Scan(&m.Type, &m.Name, &m.Value)
		if err != nil {
			tx.Rollback()
			return err
		}
		tx.Commit()
		return nil
	})
	return m, r.Retry(ctx)
}

func (s *PGStorage) UpdateMetric(ctx context.Context, m entity.Metric) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	em, err := s.GetByNameAndType(ctx, m.Name, m.Type)
	if err != nil {
		_, err = tx.ExecContext(ctx, "INSERT INTO metrics(type,name,value) VALUES ($1,$2,$3)", m.Type, m.Name, m.Value)
		if err != nil {
			tx.Rollback()
			return err
		}
	} else {
		if m.Type == entity.Counter {
			_, err = tx.ExecContext(ctx, "UPDATE metrics SET value = $1 WHERE name = $2 AND type = $3", m.Value+em.Value, m.Name, m.Type)
			if err != nil {
				tx.Rollback()
				return err
			}
		} else {
			_, err = tx.ExecContext(ctx, "UPDATE metrics SET value = $1 WHERE name = $2 AND type = $3", m.Value, m.Name, m.Type)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	tx.Commit()
	return nil
}

func (s *PGStorage) GetAllMetrics(ctx context.Context) []entity.Metric {
	allMetrics := make([]entity.Metric, 0)
	tx, err := s.DB.Begin()
	if err != nil {
		return nil
	}
	results, err := tx.QueryContext(ctx, "SELECT type,name,value FROM metrics")
	if err != nil {
		tx.Rollback()
		return nil
	}
	defer results.Close()
	for results.Next() {
		var m entity.Metric
		err = results.Scan(&m.Type, &m.Name, &m.Value)
		if err != nil {
			tx.Rollback()
			return nil
		}
		allMetrics = append(allMetrics, m)
	}
	err = results.Err()
	if err != nil {
		tx.Rollback()
		return nil
	}
	return allMetrics
}

func (s *PGStorage) ExportToFile() error {
	return nil
}

func (s *PGStorage) ImportFromFile() error {
	return nil
}

func (s *PGStorage) UpdateMetrics(ctx context.Context, mcs []entity.Metric) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	for _, m := range mcs {
		em, err := s.GetByNameAndType(ctx, m.Name, m.Type)
		if err != nil {
			_, err = tx.ExecContext(ctx, "INSERT INTO metrics(type,name,value) VALUES ($1,$2,$3)", m.Type, m.Name, m.Value)
			if err != nil {
				tx.Rollback()
				return err
			}
		} else {
			if m.Type == entity.Counter {
				_, err = tx.ExecContext(ctx, "UPDATE metrics SET value = $1 WHERE name = $2 AND type = $3", m.Value+em.Value, m.Name, m.Type)
				if err != nil {
					tx.Rollback()
					return err
				}
			} else {
				_, err = tx.ExecContext(ctx, "UPDATE metrics SET value = $1 WHERE name = $2 AND type = $3", m.Value, m.Name, m.Type)
				if err != nil {
					tx.Rollback()
					return err
				}
			}
		}
	}
	tx.Commit()
	return nil
}

func DBMigration(path string, db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://"+path,
		"postgres", driver)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	fmt.Println("Migrations applied successfully!")
	return nil
}
