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

type valueHandlerSvc interface {
	GetMetric(ctx context.Context, n string, t int64) (entity.Metric, error)
}

type ValueHandler struct {
	valueSvc valueHandlerSvc
}

func (h *ValueHandler) Handle(c echo.Context) error {
	//Var for decoding request JSON
	var reqData models.Metrics

	appHeader := c.Request().Header.Get("Content-Type")
	if appHeader != constants.ContentTypeJSON {
		c.Response().Header().Set("Content-Type", constants.ContentTypePlain)
		return c.String(http.StatusBadRequest, "Bad request")
	}

	//Reading body
	body, err := io.ReadAll(c.Request().Body)
	defer c.Request().Body.Close()
	if err != nil {
		c.Response().Header().Set("Content-Type", constants.ContentTypePlain)
		return c.String(http.StatusInternalServerError, "Server error")
	}

	//Decoding body data
	err = json.Unmarshal(body, &reqData)
	if err != nil {
		c.Response().Header().Set("Content-Type", constants.ContentTypePlain)
		return c.String(http.StatusBadRequest, "Bad request, could not unmarshal")
	}

	//Chcking id not emtpy
	if reqData.ID == "" {
		c.Response().Header().Set("Content-Type", constants.ContentTypePlain)
		return c.String(http.StatusNotFound, "Metric name not found")
	}

	if reqData.MType == "" {
		c.Response().Header().Set("Content-Type", constants.ContentTypePlain)
		return c.String(http.StatusBadRequest, "Bad request")
	}
	//Dealing request depending on type
	c.Response().Header().Set("Content-Type", constants.ContentTypeJSON)
	switch reqData.MType {
	case constants.GaugeType:
		metric, err := h.valueSvc.GetMetric(c.Request().Context(), reqData.ID, entity.Gauge)
		if err != nil {
			c.Response().Header().Set("Content-Type", constants.ContentTypePlain)
			return c.String(http.StatusNotFound, "Metric name not found")
		}
		metricValue := metric.Value
		reqData.Value = &metricValue
		return c.JSON(http.StatusOK, reqData)
	case constants.CounterType:
		metric, err := h.valueSvc.GetMetric(c.Request().Context(), reqData.ID, entity.Counter)
		if err != nil {
			c.Response().Header().Set("Content-Type", constants.ContentTypePlain)
			return c.String(http.StatusNotFound, "Metric name not found")
		}
		metricValue := int64(metric.Value)
		reqData.Delta = &metricValue
		return c.JSON(http.StatusOK, reqData)
	default:
		c.Response().Header().Set("Content-Type", constants.ContentTypePlain)
		return c.String(http.StatusBadRequest, "Bad request")
	}
}

func NewValueHandler(s valueHandlerSvc) *ValueHandler {
	return &ValueHandler{
		valueSvc: s,
	}
}
