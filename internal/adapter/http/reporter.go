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
	"os"

	"github.com/madcarpet/metrics/internal/adapter/storage"
	"github.com/madcarpet/metrics/internal/entity"
	"github.com/madcarpet/metrics/internal/models"
	"github.com/madcarpet/metrics/internal/retry"
)

func signatory(data []byte, secretKey string) string {
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write(data)
	sign := h.Sum(nil)
	return hex.EncodeToString(sign)
}

type reporter struct {
	serverAddress string
	repo          storage.Repository
	dataSign      bool
}

func NewReporter(sa string, r storage.Repository, ds bool) *reporter {
	return &reporter{serverAddress: sa, repo: r, dataSign: ds}
}

func (r *reporter) ReportMetrics() error {
	rt := retry.NewRetrier(retry.DefaultRetry, retry.Interval2s, func(ctx context.Context) error {
		var finalReqData []models.Metrics
		var body bytes.Buffer
		metrics, err := r.repo.GetAllMetrics(ctx)
		if err != nil {
			return err
		}
		for _, metric := range metrics {
			var reqData models.Metrics
			switch metric.Type {
			case entity.Gauge:
				mName := metric.Name
				mVal := metric.Value
				reqData = models.Metrics{
					ID:    mName,
					MType: "gauge",
					Value: &mVal,
				}
			case entity.Counter:
				mName := metric.Name
				metricDelta := int64(metric.Value)
				reqData = models.Metrics{
					ID:    mName,
					MType: "counter",
					Delta: &metricDelta,
				}
			}
			finalReqData = append(finalReqData, reqData)
		}
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
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Content-Encoding", "gzip")
		req.Header.Set("Accept-Encoding", "gzip")

		if r.dataSign {
			key := os.Getenv("CLIENT_SECRET_KEY")
			fmt.Println(key, "<-------------KEY in AGENT") // TODO убрать
			// sign := signatory(jsonBody, key)
			req.Header.Set("HashSHA256", "test")
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		return nil
	})
	return rt.Retry(context.Background())
}
