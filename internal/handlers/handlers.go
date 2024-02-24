package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/madcarpet/metrics/internal/storage"
)

// type for Handler with storage
type Handler struct {
	Storage storage.Repositories
}

// Constructor function for Handler
func NewHandler(s storage.Repositories) *Handler {
	return &Handler{Storage: s}
}

// Constants for metrics types
const (
	gauge storage.MetricType = 1 + iota
	counter
)

// Function fot / path
// Should be enclosed inside echo route
func (h *Handler) Root(c echo.Context) error {
	var output string
	allmetrics, err := h.Storage.GetMetricsAll()
	if err != nil {
		return c.String(http.StatusBadRequest, "Bad request")
	}
	for k, v := range allmetrics {
		output += fmt.Sprintf("%v: %v\n", k, v)
	}
	c.Response().Header().Set("Content-Type", "text/plain; charset=UTF-8")
	return c.String(http.StatusOK, output)

}

// Function for /value path
// Should be enclosed inside echo route
func (h *Handler) Value(c echo.Context) error {
	mtype := c.Param("type")
	mname := c.Param("name")
	switch mtype {
	case "gauge":
		value, err := h.Storage.GetMetric(mname, gauge)
		if err != nil {
			return c.String(http.StatusNotFound, "Metric name not found")
		}
		c.Response().Header().Set("Content-Type", "text/plain; charset=UTF-8")
		return c.String(http.StatusOK, value)
	case "counter":
		value, err := h.Storage.GetMetric(mname, counter)
		if err != nil {
			return c.String(http.StatusNotFound, "Metric name not found")
		}
		c.Response().Header().Set("Content-Type", "text/plain; charset=UTF-8")
		return c.String(http.StatusOK, value)
	}
	return c.String(http.StatusBadRequest, "Bad request")
}

// Function for /update path
// Should be enclosed inside echo route
func (h *Handler) Update(c echo.Context) error {
	mtype := c.Param("type")
	mname := c.Param("name")
	mval := c.Param("value")
	if mname == "" {
		return c.String(http.StatusNotFound, "Metric name not found")
	}
	switch mtype {
	case "gauge":
		val, err := strconv.ParseFloat(mval, 64)
		if err != nil {
			return c.String(http.StatusBadRequest, "Bad request")
		}
		h.Storage.Update(mname, val, gauge)
		c.Response().Header().Set("Content-Type", "text/plain; charset=UTF-8")
		return c.String(http.StatusOK, "Metric updated")
	case "counter":
		_, err := strconv.Atoi(mval)
		if err != nil {
			return c.String(http.StatusBadRequest, "Bad request")
		}
		val, err := strconv.ParseFloat(mval, 64)
		if err != nil {
			return c.String(http.StatusBadRequest, "Bad request")
		}
		h.Storage.Update(mname, val, counter)
		c.Response().Header().Set("Content-Type", "text/plain; charset=UTF-8")
		return c.String(http.StatusOK, "Metric updated")
	default:
		return c.String(http.StatusBadRequest, "Bad request")
	}
}
