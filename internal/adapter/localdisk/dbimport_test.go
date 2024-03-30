package localdisk

import (
	"testing"

	"github.com/madcarpet/metrics/internal/adapter/memstorage"
	"github.com/madcarpet/metrics/internal/service/metrics"
	"github.com/stretchr/testify/assert"
)

func TestDbImport(t *testing.T) {
	db := memstorage.NewMemStorage()
	updateSVC := metrics.NewUpdateMetricSvc(db)
	err := DBImport("../../../testdata/test_db.json", updateSVC)
	assert.Nil(t, err)

	metrics := db.GetAllMetrics()
	assert.Equal(t, 29, len(metrics))
	for _, m := range metrics {
		switch m.Name {
		case "Alloc":
			assert.Equal(t, float64(374096), m.Value)
		case "Pollcount":
			assert.Equal(t, float64(6), m.Value)
		case "StackSys":
			assert.Equal(t, float64(393216), m.Value)
		}
	}

}
