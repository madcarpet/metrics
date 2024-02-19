package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/madcarpet/metrics/internal/handlers"
	"github.com/madcarpet/metrics/internal/storage"
)

func main() {
	err := parseFlags()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	db := storage.NewMemStorage()
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return handlers.Root(c, db)
	})
	e.GET("/value/:type/:name", func(c echo.Context) error {
		return handlers.Value(c, db)
	})
	e.POST("/update/:type/", func(c echo.Context) error {
		c.Response().Header().Set("Content-Type", "text/plain; charset=UTF-8")
		return c.String(http.StatusNotFound, "Metric name not found")
	})
	e.POST("/update/:type/:value", func(c echo.Context) error {
		c.Response().Header().Set("Content-Type", "text/plain; charset=UTF-8")
		return c.String(http.StatusNotFound, "Metric name not found")
	})
	e.POST("/update/:type/:name/:value", func(c echo.Context) error {
		return handlers.Update(c, db)
	})
	e.Any("/*", func(c echo.Context) error {
		return c.String(http.StatusBadRequest, "Bad request")
	})
	e.Logger.Fatal(e.Start(serverAddress))
}
