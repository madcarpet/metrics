package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/madcarpet/metrics/internal/entity"
	"github.com/madcarpet/metrics/internal/models"
)

type updateHandlerSvc interface {
	UpdateMetric(m entity.Metric) error
}

type updateHandlerGetSvc interface {
	GetMetric(n string, t int64) (entity.Metric, error)
}

type UpdateHandler struct {
	updateSvc updateHandlerSvc
	getSvc    updateHandlerGetSvc
}

func (h *UpdateHandler) Handle(c echo.Context) error {
	//Var for decoding request JSON
	var updateData models.Metrics
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
	err = json.Unmarshal(body, &updateData)
	if err != nil {
		c.Response().Header().Set("Content-Type", "text/plain; charset=UTF-8")
		return c.String(http.StatusBadRequest, "Bad request")
	}

	//Chcking id not emtpy
	if updateData.ID == "" {
		c.Response().Header().Set("Content-Type", "text/plain; charset=UTF-8")
		return c.String(http.StatusNotFound, "Metric name not found")
	}

	c.Response().Header().Set("Content-Type", "application/json")

	//Dealing data, depending on type
	switch updateData.MType {
	case "gauge":
		if updateData.Value == nil {
			c.Response().Header().Set("Content-Type", "text/plain; charset=UTF-8")
			return c.String(http.StatusBadRequest, "Bad request")
		}
		metric := entity.Metric{
			Type:  entity.Gauge,
			Name:  updateData.ID,
			Value: *updateData.Value,
		}
		err = h.updateSvc.UpdateMetric(metric)
		if err != nil {
			c.Response().Header().Set("Content-Type", "text/plain; charset=UTF-8")
			return c.String(http.StatusInternalServerError, "Server error")
		}
		return c.JSON(http.StatusOK, updateData)
	case "counter":
		if updateData.Delta == nil {
			c.Response().Header().Set("Content-Type", "text/plain; charset=UTF-8")
			return c.String(http.StatusBadRequest, "Bad request")
		}
		metric := entity.Metric{
			Type:  entity.Counter,
			Name:  updateData.ID,
			Value: float64(*updateData.Delta),
		}
		err = h.updateSvc.UpdateMetric(metric)
		if err != nil {
			c.Response().Header().Set("Content-Type", "text/plain; charset=UTF-8")
			return c.String(http.StatusInternalServerError, "Server error")
		}
		currMetric, err := h.getSvc.GetMetric(updateData.ID, entity.Counter)
		if err != nil {
			c.Response().Header().Set("Content-Type", "text/plain; charset=UTF-8")
			return c.String(http.StatusInternalServerError, "Server error")
		}
		newValue := currMetric.Value
		updateData.Value = &newValue
		return c.JSON(http.StatusOK, updateData)

	default:
		c.Response().Header().Set("Content-Type", "text/plain; charset=UTF-8")
		return c.String(http.StatusBadRequest, "Bad request")
	}
}

func NewUpdateHandler(us updateHandlerSvc, gs updateHandlerGetSvc) *UpdateHandler {
	return &UpdateHandler{
		updateSvc: us,
		getSvc:    gs,
	}
}
