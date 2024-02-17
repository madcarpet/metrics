package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/madcarpet/metrics/internal/storage"
)

// Constants gor metrics types
const (
	gauge storage.MetricType = 1 + iota
	counter
)

// Interface to work with MemStorage
type Repositories interface {
	Update(k string, v float64, t storage.MetricType)
	GetMetric(k string, t storage.MetricType) (string, error)
	GetMetricsAll() (map[string]string, error)
}

func Root(c echo.Context, s Repositories) error {
	var output string
	allmetrics, err := s.GetMetricsAll()
	if err != nil {
		return c.String(http.StatusBadRequest, "Bad request")
	}
	for k, v := range allmetrics {
		output += fmt.Sprintf("%v: %v\n", k, v)
	}
	c.Response().Header().Set("Content-Type", "text/plain; charset=UTF-8")
	return c.String(http.StatusOK, output)

}

// Handler fot /value path
// Should be enclosed
func Value(c echo.Context, s Repositories) error {
	mtype := c.Param("type")
	mname := c.Param("name")
	switch mtype {
	case "gauge":
		value, err := s.GetMetric(mname, gauge)
		if err != nil {
			return c.String(http.StatusNotFound, "Metric name not found")
		}
		c.Response().Header().Set("Content-Type", "text/plain; charset=UTF-8")
		return c.String(http.StatusOK, value)
	case "counter":
		value, err := s.GetMetric(mname, counter)
		if err != nil {
			return c.String(http.StatusNotFound, "Metric name not found")
		}
		c.Response().Header().Set("Content-Type", "text/plain; charset=UTF-8")
		return c.String(http.StatusOK, value)
	}
	return c.String(http.StatusBadRequest, "Bad request")
}

// Handler for /update path
// Should be enclosed
func Update(c echo.Context, s Repositories) error {
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
		s.Update(mname, val, gauge)
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
		s.Update(mname, val, counter)
		c.Response().Header().Set("Content-Type", "text/plain; charset=UTF-8")
		return c.String(http.StatusOK, "Metric updated")
	default:
		return c.String(http.StatusBadRequest, "Bad request")
	}
}
