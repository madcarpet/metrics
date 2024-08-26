package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/madcarpet/metrics/internal/entity"
	"github.com/madcarpet/metrics/internal/models"
)

type updatesHandlerSvc interface {
	UpdateMetrics(ctx context.Context, m []entity.Metric) error
}

type updatesHandlerGetSvc interface {
	GetMetric(ctx context.Context, n string, t int64) (entity.Metric, error)
}

type UpdatesHandler struct {
	updatesSvc updatesHandlerSvc
	getSvc     updatesHandlerGetSvc
}

func NewUpdatesHandler(us updatesHandlerSvc, gs updatesHandlerGetSvc) *UpdatesHandler {
	return &UpdatesHandler{
		updatesSvc: us,
		getSvc:     gs,
	}
}

func (h *UpdatesHandler) Handle(c echo.Context) error {
	//var for checking existance
	existGuges := make(map[string]entity.Metric)
	existCounters := make(map[string]entity.Metric)
	//Var for decoding request JSON
	var gotData, respData []models.Metrics
	var updateData []entity.Metric

	//Checking Content-Type header
	appHeader := c.Request().Header.Get("Content-Type")
	if appHeader != "application/json" {
		c.Response().Header().Set("Content-Type", "text/plain; charset=UTF-8")
		return c.String(http.StatusBadRequest, "Bad request")
	}
	//Reading body
	body, err := io.ReadAll(c.Request().Body)
	defer c.Request().Body.Close()
	if err != nil {
		c.Response().Header().Set("Content-Type", "text/plain; charset=UTF-8")
		return c.String(http.StatusInternalServerError, "Server error")
	}

	//Decoding body data
	err = json.Unmarshal(body, &gotData)
	if err != nil {
		c.Response().Header().Set("Content-Type", "text/plain; charset=UTF-8")
		return c.String(http.StatusBadRequest, "Bad request")
	}

	//Chcking metrics is valid
	fmt.Println(gotData)
	for _, m := range gotData {
		if m.ID == "" {
			c.Response().Header().Set("Content-Type", "text/plain; charset=UTF-8")
			return c.String(http.StatusNotFound, "Metric name not found for some metrics")
		}
		//Dealing data, depending on type
		switch m.MType {
		case "gauge":
			if m.Value == nil {
				c.Response().Header().Set("Content-Type", "text/plain; charset=UTF-8")
				return c.String(http.StatusBadRequest, "Bad request")
			}
			metric := entity.Metric{
				Type:  entity.Gauge,
				Name:  m.ID,
				Value: *m.Value,
			}
			existGuges[metric.Name] = metric
		case "counter":
			if m.Delta == nil {
				c.Response().Header().Set("Content-Type", "text/plain; charset=UTF-8")
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
			c.Response().Header().Set("Content-Type", "text/plain; charset=UTF-8")
			return c.String(http.StatusBadRequest, "Bad request")
		}
	}
	fmt.Println(existGuges)
	for _, metric := range existGuges {
		updateData = append(updateData, metric)
		respMetric := models.Metrics{
			ID:    metric.Name,
			MType: "gauge",
			Value: &metric.Value,
		}
		respData = append(respData, respMetric)
	}
	for _, metric := range existCounters {
		updateData = append(updateData, metric)
		newValue := int64(metric.Value)
		respMetric := models.Metrics{
			ID:    metric.Name,
			MType: "counter",
			Delta: &newValue,
		}
		respData = append(respData, respMetric)
	}
	c.Response().Header().Set("Content-Type", "application/json")
	err = h.updatesSvc.UpdateMetrics(c.Request().Context(), updateData)
	if err != nil {
		c.Response().Header().Set("Content-Type", "text/plain; charset=UTF-8")
		return c.String(http.StatusInternalServerError, "Server error")
	}
	return c.JSON(http.StatusOK, respData)
}
