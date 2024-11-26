package handlers

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
)

// pingHandelerSvc interface for ping service.
type pingHandelerSvc interface {
	Ping(ctx context.Context) error
}

// PingHandler struct for ping handler keeps ping service.
type PingHandler struct {
	pingSvc pingHandelerSvc
}

// NewPingHandler creates a new ping handler.
func NewPingHandler(ctx context.Context, s pingHandelerSvc) *PingHandler {
	return &PingHandler{pingSvc: s}
}

// Handle handles http request.
func (h *PingHandler) Handle(c echo.Context) error {
	err := h.pingSvc.Ping(c.Request().Context())
	if err != nil {
		return c.String(http.StatusInternalServerError, "500 Internal Server Error")
	}
	return c.String(http.StatusOK, "Connection to database established")
}
