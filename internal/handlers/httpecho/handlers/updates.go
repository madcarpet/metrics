package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/madcarpet/metrics/internal/constants"
	"github.com/madcarpet/metrics/internal/entity"
	"github.com/madcarpet/metrics/internal/models"
)

// updatesHandlerSvc interface with metod to update metric data.
type updatesHandlerSvc interface {
	UpdateMetrics(ctx context.Context, m []entity.Metric) error
}

// updatesHandlerGetSvc interface with metod to get metric data.
type updatesHandlerGetSvc interface {
	GetMetric(ctx context.Context, n string, t int64) (entity.Metric, error)
}

// UpdatesHandler structure for updating metrics handler keeps update metric and get metric services.
type UpdatesHandler struct {
	updatesSvc updatesHandlerSvc
	getSvc     updatesHandlerGetSvc
}

// NewUpdatesHandler creates a new updatesHandler.
func NewUpdatesHandler(us updatesHandlerSvc, gs updatesHandlerGetSvc) *UpdatesHandler {
	return &UpdatesHandler{
		updatesSvc: us,
		getSvc:     gs,
	}
}

// Handle handles http request.
func (h *UpdatesHandler) Handle(c echo.Context) error {

	// Var for decoding request JSON.
	existGuges := make(map[string]entity.Metric)
	existCounters := make(map[string]entity.Metric)
	var gotData []models.Metrics

	// Checking Content-Type header.
	appHeader := c.Request().Header.Get("Content-Type")
	if appHeader != constants.ContentTypeJSON {
		c.Response().Header().Set("Content-Type", constants.ContentTypePlain)
		return c.String(http.StatusBadRequest, "Bad request")
	}
	// Reading body.
	body, err := io.ReadAll(c.Request().Body)
	defer c.Request().Body.Close()
	if err != nil {
		c.Response().Header().Set("Content-Type", constants.ContentTypePlain)
		return c.String(http.StatusInternalServerError, "Server error")
	}

	// Decoding body data.
	err = json.Unmarshal(body, &gotData)
	if err != nil {
		c.Response().Header().Set("Content-Type", constants.ContentTypePlain)
		return c.String(http.StatusBadRequest, "Bad request")
	}

	// Var for checking existance.
	respData := make([]models.Metrics, 0, len(gotData))

	// Checking metrics is valid.
	for _, m := range gotData {
		if m.ID == "" {
			c.Response().Header().Set("Content-Type", constants.ContentTypePlain)
			return c.String(http.StatusNotFound, "Metric name not found for some metrics")
		}
		// Dealing data, depending on type.
		switch m.MType {
		case constants.GaugeType:
			if m.Value == nil {
				c.Response().Header().Set("Content-Type", constants.ContentTypePlain)
				return c.String(http.StatusBadRequest, "Bad request")
			}
			metric := entity.Metric{
				Type:  entity.Gauge,
				Name:  m.ID,
				Value: *m.Value,
			}
			existGuges[metric.Name] = metric
		case constants.CounterType:
			if m.Delta == nil {
				c.Response().Header().Set("Content-Type", constants.ContentTypePlain)
				return c.String(http.StatusBadRequest, "Bad request")
			}
			metric := entity.Metric{
				Type:  entity.Counter,
				Name:  m.ID,
				Value: float64(*m.Delta),
			}
			if existCounter, exists := existCounters[metric.Name]; exists {
				existCounter.Value += metric.Value
				existCounters[metric.Name] = existCounter
			} else {
				existCounters[metric.Name] = metric
			}
		default:
			c.Response().Header().Set("Content-Type", constants.ContentTypePlain)
			return c.String(http.StatusBadRequest, "Bad request")
		}
	}
	updateData := make([]entity.Metric, 0, len(gotData))
	for _, metric := range existGuges {
		updateData = append(updateData, metric)
		respMetric := models.Metrics{
			ID:    metric.Name,
			MType: constants.GaugeType,
			Value: &metric.Value,
		}
		respData = append(respData, respMetric)
	}
	for _, metric := range existCounters {
		updateData = append(updateData, metric)
		newValue := int64(metric.Value)
		respMetric := models.Metrics{
			ID:    metric.Name,
			MType: constants.CounterType,
			Delta: &newValue,
		}
		respData = append(respData, respMetric)
	}
	c.Response().Header().Set("Content-Type", constants.ContentTypeJSON)
	err = h.updatesSvc.UpdateMetrics(c.Request().Context(), updateData)
	if err != nil {
		c.Response().Header().Set("Content-Type", constants.ContentTypePlain)
		return c.String(http.StatusInternalServerError, "Server error")
	}
	return c.JSON(http.StatusOK, respData)
}
