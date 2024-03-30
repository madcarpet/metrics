package localdisk

import (
	"encoding/json"
	"os"

	"github.com/madcarpet/metrics/internal/entity"
)

type repository interface {
	GetAllMetrics() []entity.Metric
}

type dbSaver struct {
	path string
	repo repository
}

func NewSaver(dp string, r repository) *dbSaver {
	return &dbSaver{path: dp, repo: r}
}

func (s *dbSaver) WriteDb() error {
	metrics := s.repo.GetAllMetrics()
	if len(metrics) > 0 {
		flag := os.O_WRONLY | os.O_CREATE
		df, err := os.OpenFile(s.path, flag, 0644)
		if err != nil {
			return err
		}
		data, err := json.Marshal(metrics)
		if err != nil {
			return err
		}
		_, err = df.Write(data)
		if err != nil {
			return err
		}
		if err = df.Close(); err != nil {
			return err
		}
	}
	return nil
}
