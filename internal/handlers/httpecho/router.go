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
	updateHandler := handlers.NewUpdateHandler(updateSvc, valueSvc)

	valueURLHandler := handlers.NewValueURLHandler(valueSvc)
	updateURLHandler := handlers.NewUpdateURLHandler(updateSvc)

	// Root handling
	e.GET("/", rootHandler.Handle, middlewares.ReqRespWithLogging, middlewares.GzipCompression)
	// JSON requests handling
	e.POST("/value/", valueHandler.Handle, middlewares.ReqRespWithLogging, middlewares.GzipCompression)
	e.POST("/update/", updateHandler.Handle, middlewares.ReqRespWithLogging, middlewares.GzipCompression)
	// Requests via URL handling
	e.GET("/value/:type/:name", valueURLHandler.Handle, middlewares.ReqRespWithLogging)
	e.POST("/update/:type/", func(c echo.Context) error {
		c.Response().Header().Set("Content-Type", "text/plain; charset=UTF-8")
		return c.String(http.StatusNotFound, "Metric name not found")
	}, middlewares.ReqRespWithLogging)
	e.POST("/update/:type/:value", func(c echo.Context) error {
		c.Response().Header().Set("Content-Type", "text/plain; charset=UTF-8")
		return c.String(http.StatusNotFound, "Metric name not found")
	}, middlewares.ReqRespWithLogging)
	e.POST("/update/:type/:name/:value", updateURLHandler.Handle, middlewares.ReqRespWithLogging, middlewares.GzipCompression)
	// Any handling
	e.Any("/*", func(c echo.Context) error {
		c.Response().Header().Set("Content-Type", "text/plain; charset=UTF-8")
		return c.String(http.StatusBadRequest, "Bad request")
	}, middlewares.ReqRespWithLogging)
}
