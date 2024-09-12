package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/madcarpet/metrics/internal/constants"
	"github.com/madcarpet/metrics/internal/entity"
)

type valueURLHandlerSvc interface {
	GetMetric(ctx context.Context, n string, t int64) (entity.Metric, error)
}

type ValueURLHandler struct {
	valueSvc valueHandlerSvc
}

func (h *ValueURLHandler) Handle(c echo.Context) error {
	mType := c.Param("type")
	mName := c.Param("name")
	switch mType {
	case constants.GaugeType:
		metric, err := h.valueSvc.GetMetric(c.Request().Context(), mName, entity.Gauge)
		if err != nil {
			return c.String(http.StatusNotFound, "Metric name not found")
		}
		c.Response().Header().Set("Content-Type", constants.ContentTypePlain)
		return c.String(http.StatusOK, fmt.Sprintf("%v", metric.Value))
	case constants.CounterType:
		metric, err := h.valueSvc.GetMetric(c.Request().Context(), mName, entity.Counter)
		if err != nil {
			return c.String(http.StatusNotFound, "Metric name not found")
		}
		c.Response().Header().Set("Content-Type", constants.ContentTypePlain)
		return c.String(http.StatusOK, fmt.Sprintf("%v", metric.Value))
	}
	return c.String(http.StatusBadRequest, "Bad request")

}

func NewValueURLHandler(s valueURLHandlerSvc) *ValueURLHandler {
	return &ValueURLHandler{
		valueSvc: s,
	}
}
