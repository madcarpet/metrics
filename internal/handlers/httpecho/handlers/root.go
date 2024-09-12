package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/madcarpet/metrics/internal/constants"
	"github.com/madcarpet/metrics/internal/entity"
)

type rootHandlerSvc interface {
	GetAllMetrics(ctx context.Context) ([]entity.Metric, error)
}

type RootHandler struct {
	rootSvc rootHandlerSvc
}

func (r *RootHandler) Handle(c echo.Context) error {
	var output string
	allMetrics, err := r.rootSvc.GetAllMetrics(c.Request().Context())
	if err != nil {
		c.Response().Header().Set("Content-Type", constants.ContentTypePlain)
		return c.String(http.StatusInternalServerError, "Server error")
	}
	for _, m := range allMetrics {
		output += fmt.Sprintf("%v: %v\n", m.Name, m.Value)
	}
	c.Response().Header().Set("Content-Type", "text/html")
	return c.String(http.StatusOK, output)
}

func NewRootHandler(s rootHandlerSvc) *RootHandler {
	return &RootHandler{
		rootSvc: s,
	}
}
