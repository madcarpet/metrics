package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/madcarpet/metrics/internal/entity"
)

type valueHandlerSvc interface {
	GetMetric(n string, t int64) (entity.Metric, error)
}

type ValueHandler struct {
	valueSvc valueHandlerSvc
}

func (h *ValueHandler) Handle(c echo.Context) error {
	mType := c.Param("type")
	mName := c.Param("name")
	switch mType {
	case "gauge":
		metric, err := h.valueSvc.GetMetric(mName, entity.Gauge)
		if err != nil {
			return c.String(http.StatusNotFound, "Metric name not found")
		}
		c.Response().Header().Set("Content-Type", "text/plain; charset=UTF-8")
		return c.String(http.StatusOK, fmt.Sprintf("%v", metric.Value))
	case "counter":
		metric, err := h.valueSvc.GetMetric(mName, entity.Counter)
		if err != nil {
			return c.String(http.StatusNotFound, "Metric name not found")
		}
		c.Response().Header().Set("Content-Type", "text/plain; charset=UTF-8")
		return c.String(http.StatusOK, fmt.Sprintf("%v", metric.Value))
	}
	return c.String(http.StatusBadRequest, "Bad request")

}

func NewValueHandler(s valueHandlerSvc) *ValueHandler {
	return &ValueHandler{
		valueSvc: s,
	}
}
