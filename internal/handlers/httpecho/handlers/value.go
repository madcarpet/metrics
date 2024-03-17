package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/madcarpet/metrics/internal/entity"
	"github.com/madcarpet/metrics/internal/models"
)

type valueHandlerSvc interface {
	GetMetric(n string, t int64) (entity.Metric, error)
}

type ValueHandler struct {
	valueSvc valueHandlerSvc
}

func (h *ValueHandler) Handle(c echo.Context) error {
	//Var for decoding request JSON
	var reqData models.Metrics

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
	err = json.Unmarshal(body, &reqData)
	if err != nil {
		c.Response().Header().Set("Content-Type", "text/plain; charset=UTF-8")
		return c.String(http.StatusBadRequest, "Bad request")
	}

	//Chcking id not emtpy
	if reqData.ID == "" {
		c.Response().Header().Set("Content-Type", "text/plain; charset=UTF-8")
		return c.String(http.StatusNotFound, "Metric name not found")
	}

	if reqData.MType == "" {
		c.Response().Header().Set("Content-Type", "text/plain; charset=UTF-8")
		return c.String(http.StatusBadRequest, "Bad request")
	}
	//Dealing request depending on type
	c.Response().Header().Set("Content-Type", "application/json")
	switch reqData.MType {
	case "gauge":
		metric, err := h.valueSvc.GetMetric(reqData.ID, entity.Gauge)
		if err != nil {
			c.Response().Header().Set("Content-Type", "text/plain; charset=UTF-8")
			return c.String(http.StatusNotFound, "Metric name not found")
		}
		metricValue := metric.Value
		reqData.Value = &metricValue
		return c.JSON(http.StatusOK, reqData)
	case "counter":
		metric, err := h.valueSvc.GetMetric(reqData.ID, entity.Counter)
		if err != nil {
			c.Response().Header().Set("Content-Type", "text/plain; charset=UTF-8")
			return c.String(http.StatusNotFound, "Metric name not found")
		}
		metricValue := metric.Value
		reqData.Value = &metricValue
		return c.JSON(http.StatusOK, reqData)
	default:
		c.Response().Header().Set("Content-Type", "text/plain; charset=UTF-8")
		return c.String(http.StatusBadRequest, "Bad request")
	}
}

func NewValueHandler(s valueHandlerSvc) *ValueHandler {
	return &ValueHandler{
		valueSvc: s,
	}
}
