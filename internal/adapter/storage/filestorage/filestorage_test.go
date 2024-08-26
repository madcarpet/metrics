package filestorage

import (
	"context"
	"os"
	"testing"

	"github.com/madcarpet/metrics/internal/entity"
	"github.com/stretchr/testify/assert"
)

func TestFileStorage(t *testing.T) {
	testFile := "test.json"
	storage, err := NewFileStorage(testFile, false)
	assert.Nil(t, err)
	_, err = os.Stat(testFile)
	assert.Nil(t, err)

	updateMetrics := []entity.Metric{
		{Name: "TestGauge1", Type: entity.Gauge, Value: 1.111111},
		{Name: "TestCounter1", Type: entity.Counter, Value: 99},
	}

	checkMetrics := []string{`{"type":1,"name":"TestGauge1","value":1.111111}`, `{"type":2,"name":"TestCounter1","value":99}`}

	for _, metric := range updateMetrics {
		storage.UpdateMetric(context.TODO(), metric)
	}

	err = storage.ExportToFile()
	assert.Nil(t, err)

	data, err := os.ReadFile(testFile)
	assert.Nil(t, err)

	for _, val := range checkMetrics {
		assert.Contains(t, string(data), val)
	}

	storageImport, err := NewFileStorage(testFile, false)
	assert.Nil(t, err)
	err = storageImport.ImportFromFile()
	assert.Nil(t, err)

	for _, met := range updateMetrics {
		metricImported, err := storage.GetByNameAndType(context.TODO(), met.Name, met.Type)
		assert.Nil(t, err)
		assert.Equal(t, metricImported, met)
	}
	storageImport.file.Close()

	err = storage.Close()
	assert.Nil(t, err)

	err = os.Remove(testFile)
	assert.Nil(t, err)

}
