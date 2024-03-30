package localdisk

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/madcarpet/metrics/internal/adapter/memstorage"
	"github.com/madcarpet/metrics/internal/entity"
	"github.com/stretchr/testify/assert"
)

func TestDbSaver(t *testing.T) {
	db := memstorage.NewMemStorage()
	testMetrics := []entity.Metric{
		{Name: "Alloc", Type: entity.Gauge, Value: 2.678488e+06},
		{Name: "BuckHashSys", Type: entity.Gauge, Value: 6160},
		{Name: "GCCPUFraction", Type: entity.Gauge, Value: 0},
		{Name: "RandomValue", Type: entity.Gauge, Value: 0.9233588813342314},
		{Name: "PollCount", Type: entity.Counter, Value: 991112111111},
	}
	for _, metric := range testMetrics {
		db.UpdateMetric(metric)
	}
	path := "../../../testdata/test_db_saver.json"

	testDaver := NewSaver(path, db)
	testDaver.WriteDB()

	_, err := os.Stat(path)
	if err != nil {
		t.Errorf("DB file was not created")
	}

	testFile, _ := os.OpenFile(path, os.O_RDONLY, 0655)

	var decodedData []entity.Metric
	decoder := json.NewDecoder(testFile)
	decoder.Decode(&decodedData)
	testFile.Close()

	for _, m := range decodedData {
		switch m.Name {
		case "Alloc":
			assert.Equal(t, float64(2.678488e+06), m.Value)
		case "PollCount":
			assert.Equal(t, float64(991112111111), m.Value)
		case "BuckHashSys":
			assert.Equal(t, float64(6160), m.Value)
		}
	}

	err = os.Remove(path)
	if err != nil {
		t.Errorf("Test DB file was not deleted properly: %v", err)
	}

}
