package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/madcarpet/metrics/internal/constants"
	"github.com/madcarpet/metrics/internal/entity"
)

// updateURLHandlerSvc interface for updating metric service.
type updateURLHandlerSvc interface {
	UpdateMetric(ctx context.Context, m entity.Metric) error
}

// UpdateURLHandler struct for updating metric with url handler, keeps update service.
type UpdateURLHandler struct {
	updateSvc updateHandlerSvc
}

// NewUpdateURLHandler creates a new update URL handler.
func NewUpdateURLHandler(s updateURLHandlerSvc) *UpdateURLHandler {
	return &UpdateURLHandler{
		updateSvc: s,
	}
}

// Handle handles http request.
func (h *UpdateURLHandler) Handle(c echo.Context) error {
	mType := c.Param("type")
	mName := c.Param("name")
	mVal := c.Param("value")
	c.Response().Header().Set("Content-Type", constants.ContentTypePlain)
	if mName == "" {
		return c.String(http.StatusNotFound, "Metric name not found")
	}

	switch mType {
	case constants.GaugeType:
		val, err := strconv.ParseFloat(mVal, 64)
		if err != nil {
			return c.String(http.StatusBadRequest, "Bad request")
		}
		metric := entity.Metric{
			Type:  entity.Gauge,
			Name:  mName,
			Value: val,
		}
		err = h.updateSvc.UpdateMetric(c.Request().Context(), metric)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Server error")
		}
		return c.String(http.StatusOK, "Metric updated")
	case constants.CounterType:
		_, err := strconv.Atoi(mVal)
		if err != nil {
			return c.String(http.StatusBadRequest, "Bad request")
		}
		val, err := strconv.ParseFloat(mVal, 64)
		if err != nil {
			return c.String(http.StatusBadRequest, "Bad request")
		}
		metric := entity.Metric{
			Type:  entity.Counter,
			Name:  mName,
			Value: val,
		}
		err = h.updateSvc.UpdateMetric(c.Request().Context(), metric)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Server error")
		}
		return c.String(http.StatusOK, "Metric updated")
	default:
		return c.String(http.StatusBadRequest, "Bad request")
	}
}
