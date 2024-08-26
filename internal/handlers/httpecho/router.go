package httpecho

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/madcarpet/metrics/internal/entity"
	"github.com/madcarpet/metrics/internal/handlers/httpecho/handlers"
	"github.com/madcarpet/metrics/internal/handlers/httpecho/middlewares"
)

type rootHandlerSvc interface {
	GetAllMetrics() []entity.Metric
}
type updateHandlerSvc interface {
	UpdateMetric(m entity.Metric) error
}
type valueHandlerSvc interface {
	GetMetric(n string, t int64) (entity.Metric, error)
}

func SetupRouter(
	e *echo.Echo,
	rootSvc rootHandlerSvc,
	valueSvc valueHandlerSvc,
	updateSvc updateHandlerSvc,
) {
	rootHandler := handlers.NewRootHandler(rootSvc)
	valueHandler := handlers.NewValueHandler(valueSvc)
	updateHandler := handlers.NewUpdateHandler(updateSvc)

	e.GET("/", rootHandler.Handle, middlewares.ReqRespWithLogging)
	e.GET("/value/:type/:name", valueHandler.Handle, middlewares.ReqRespWithLogging)
	e.POST("/update/:type/", func(c echo.Context) error {
		c.Response().Header().Set("Content-Type", "text/plain; charset=UTF-8")
		return c.String(http.StatusNotFound, "Metric name not found")
	}, middlewares.ReqRespWithLogging)
	e.POST("/update/:type/:value", func(c echo.Context) error {
		c.Response().Header().Set("Content-Type", "text/plain; charset=UTF-8")
		return c.String(http.StatusNotFound, "Metric name not found")
	}, middlewares.ReqRespWithLogging)
	e.POST("/update/:type/:name/:value", updateHandler.Handle, middlewares.ReqRespWithLogging)
	e.Any("/*", func(c echo.Context) error {
		c.Response().Header().Set("Content-Type", "text/plain; charset=UTF-8")
		return c.String(http.StatusBadRequest, "Bad request")
	}, middlewares.ReqRespWithLogging)
}
