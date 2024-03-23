package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/madcarpet/metrics/internal/entity"
)

type updateURLHandlerSvc interface {
	UpdateMetric(m entity.Metric) error
}

type UpdateURLHandler struct {
	updateSvc updateHandlerSvc
}

func (h *UpdateURLHandler) Handle(c echo.Context) error {
	mType := c.Param("type")
	mName := c.Param("name")
	mVal := c.Param("value")
	c.Response().Header().Set("Content-Type", "text/plain; charset=UTF-8")
	if mName == "" {
		return c.String(http.StatusNotFound, "Metric name not found")
	}

	switch mType {
	case "gauge":
		val, err := strconv.ParseFloat(mVal, 64)
		if err != nil {
			return c.String(http.StatusBadRequest, "Bad request")
		}
		metric := entity.Metric{
			Type:  entity.Gauge,
			Name:  mName,
			Value: val,
		}
		err = h.updateSvc.UpdateMetric(metric)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Server error")
		}
		return c.String(http.StatusOK, "Metric updated")
	case "counter":
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
		err = h.updateSvc.UpdateMetric(metric)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Server error")
		}
		return c.String(http.StatusOK, "Metric updated")
	default:
		return c.String(http.StatusBadRequest, "Bad request")
	}
}

func NewUpdateURLHandler(s updateURLHandlerSvc) *UpdateURLHandler {
	return &UpdateURLHandler{
		updateSvc: s,
	}
}
