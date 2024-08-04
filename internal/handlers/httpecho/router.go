package httpecho

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/madcarpet/metrics/internal/entity"
	"github.com/madcarpet/metrics/internal/handlers/httpecho/handlers"
	"github.com/madcarpet/metrics/internal/handlers/httpecho/middlewares"
)

type rootHandlerSvc interface {
	GetAllMetrics(ctx context.Context) ([]entity.Metric, error)
}
type updateHandlerSvc interface {
	UpdateMetric(ctx context.Context, m entity.Metric) error
}
type valueHandlerSvc interface {
	GetMetric(ctx context.Context, n string, t int64) (entity.Metric, error)
}

type updatesHandlerSvc interface {
	UpdateMetrics(ctx context.Context, mcs []entity.Metric) error
}

type pingHandelerSvc interface {
	Ping(ctx context.Context) error
}

func SetupRouter(
	e *echo.Echo,
	rootSvc rootHandlerSvc,
	valueSvc valueHandlerSvc,
	updateSvc updateHandlerSvc,
	updatesSvc updatesHandlerSvc,
	pingSvc pingHandelerSvc,
	dataSign bool,
) {
	rootHandler := handlers.NewRootHandler(rootSvc)
	valueHandler := handlers.NewValueHandler(valueSvc)
	updateHandler := handlers.NewUpdateHandler(updateSvc, valueSvc)
	updatesHandler := handlers.NewUpdatesHandler(updatesSvc, valueSvc)
	pingHandler := handlers.NewPingHandler(context.Background(), pingSvc)

	valueURLHandler := handlers.NewValueURLHandler(valueSvc)
	updateURLHandler := handlers.NewUpdateURLHandler(updateSvc)

	var mwList []echo.MiddlewareFunc

	if dataSign {
		mwList = []echo.MiddlewareFunc{middlewares.ReqRespWithLogging, middlewares.GzipCompression, middlewares.SignData}
	} else {
		mwList = []echo.MiddlewareFunc{middlewares.ReqRespWithLogging, middlewares.GzipCompression}
	}

	// Root handling
	e.GET("/", rootHandler.Handle, mwList...)
	// JSON requests handling
	e.POST("/value/", valueHandler.Handle, mwList...)
	e.POST("/update/", updateHandler.Handle, mwList...)
	e.POST("/updates/", updatesHandler.Handle, mwList...)
	// Requests via URL handling
	e.GET("/value/:type/:name", valueURLHandler.Handle, middlewares.ReqRespWithLogging)
	//Ping DB
	e.GET("/ping", pingHandler.Handle, mwList...)
	e.POST("/update/:type/", func(c echo.Context) error {
		c.Response().Header().Set("Content-Type", "text/plain; charset=UTF-8")
		return c.String(http.StatusNotFound, "Metric name not found")
	}, middlewares.ReqRespWithLogging)
	e.POST("/update/:type/:value", func(c echo.Context) error {
		c.Response().Header().Set("Content-Type", "text/plain; charset=UTF-8")
		return c.String(http.StatusNotFound, "Metric name not found")
	}, middlewares.ReqRespWithLogging)
	e.POST("/update/:type/:name/:value", updateURLHandler.Handle, mwList...)
	// Any handling
	e.Any("/*", func(c echo.Context) error {
		c.Response().Header().Set("Content-Type", "text/plain; charset=UTF-8")
		return c.String(http.StatusBadRequest, "Bad request")
	}, middlewares.ReqRespWithLogging)
}
