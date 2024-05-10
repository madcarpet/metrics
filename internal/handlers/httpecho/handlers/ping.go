package handlers

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
)

type pingHandelerSvc interface {
	Ping(ctx context.Context) error
}

type PingHandler struct {
	pingSvc pingHandelerSvc
	ctx     context.Context
}

func (h *PingHandler) Handle(c echo.Context) error {
	err := h.pingSvc.Ping(h.ctx)
	if err != nil {
		return c.String(http.StatusInternalServerError, "500 Internal Server Error")
	}
	return c.String(http.StatusOK, "Connection to database established")
}

func NewPingHandler(ctx context.Context, s pingHandelerSvc) *PingHandler {
	return &PingHandler{pingSvc: s, ctx: ctx}
}
