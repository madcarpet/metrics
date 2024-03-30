package localdisk

import (
	"encoding/json"
	"os"

	"github.com/madcarpet/metrics/internal/entity"
	"github.com/madcarpet/metrics/internal/service/metrics"
)

func DBImport(path string, us *metrics.UpdateMetricSvc) error {
	if path != "" {
		var metrics []entity.Metric
		sf, err := os.OpenFile(path, os.O_RDONLY, 0655)
		if err != nil {
			return err
		}

		fileInfo, err := sf.Stat()
		if err != nil {
			return err
		}

		fData := make([]byte, fileInfo.Size())

		_, err = sf.Read(fData)
		if err != nil {
			return err
		}
		if err = sf.Close(); err != nil {
			return err
		}
		err = json.Unmarshal(fData, &metrics)
		if err != nil {
			return err
		}
		if len(metrics) > 0 {
			for _, m := range metrics {
				err = us.UpdateMetric(m)
				if err != nil {
					return err
				}
			}
		}

	}
	return nil
}
