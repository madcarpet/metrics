package http

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/madcarpet/metrics/internal/constants"
	"github.com/madcarpet/metrics/internal/entity"
	"github.com/madcarpet/metrics/internal/models"
	"github.com/madcarpet/metrics/internal/retry"
)

// signatory - function to get hash from data with signature.
func signatory(data []byte, secretKey string) string {
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write(data)
	sign := h.Sum(nil)
	return hex.EncodeToString(sign)
}

// reporter - structure for metric reporter.
type reporter struct {
	serverAddress string
	dataSign      bool
	key           string
}

// NewReporter - metric reporter constructor.
func NewReporter(sa string, ds bool, sk string) *reporter {
	return &reporter{serverAddress: sa, dataSign: ds, key: sk}
}

// ReportMetrics - function to report metrics to server by HTTP.
func (r *reporter) ReportMetrics(ctx context.Context, metrics []entity.Metric) error {
	// retry construction if server not available (network timeouts).
	rt := retry.NewRetrier(retry.DefaultRetry, retry.Interval2s, func(ctx context.Context) error {
		var finalReqData []models.Metrics
		var body bytes.Buffer
		for _, metric := range metrics {
			var reqData models.Metrics
			switch metric.Type {
			case entity.Gauge:
				mName := metric.Name
				mVal := metric.Value
				reqData = models.Metrics{
					ID:    mName,
					MType: constants.GaugeType,
					Value: &mVal,
				}
			case entity.Counter:
				mName := metric.Name
				metricDelta := int64(metric.Value)
				reqData = models.Metrics{
					ID:    mName,
					MType: constants.CounterType,
					Delta: &metricDelta,
				}
			}
			finalReqData = append(finalReqData, reqData)
		}
		// serialize data for request.
		jsonBody, err := json.Marshal(&finalReqData)
		if err != nil {
			return fmt.Errorf("report encoding error: %s", err)
		}

		gzBodyWriter := gzip.NewWriter(&body)
		_, err = gzBodyWriter.Write(jsonBody)
		if err != nil {
			return fmt.Errorf("json body compression to buffer error: %s", err)
		}
		if err := gzBodyWriter.Close(); err != nil {
			return fmt.Errorf("gzip writer closing error: %s", err)
		}
		url := fmt.Sprintf("http://%v/updates/", r.serverAddress)
		req, err := http.NewRequest(http.MethodPost, url, &body)
		if err != nil {
			return fmt.Errorf("request formation error: %s", err)
		}
		req.Header.Set("Content-Type", constants.ContentTypeJSON)
		req.Header.Set("Content-Encoding", "gzip")
		req.Header.Set("Accept-Encoding", "gzip")

		if r.dataSign {
			sign := signatory(jsonBody, r.key)
			req.Header.Set("HashSHA256", sign)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		return nil
	})
	return rt.Retry(ctx)
}
