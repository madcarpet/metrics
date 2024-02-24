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
	// initialize new storage
	db := storage.NewMemStorage()
	// initialize new echo instance
	e := echo.New()
	// initialize handler
	h := handlers.NewHandler(db)
	// routing
	e.GET("/", h.Root)
	e.GET("/value/:type/:name", h.Value)
	e.POST("/update/:type/", func(c echo.Context) error {
		c.Response().Header().Set("Content-Type", "text/plain; charset=UTF-8")
		return c.String(http.StatusNotFound, "Metric name not found")
	})
	e.POST("/update/:type/:value", func(c echo.Context) error {
		c.Response().Header().Set("Content-Type", "text/plain; charset=UTF-8")
		return c.String(http.StatusNotFound, "Metric name not found")
	})
	e.POST("/update/:type/:name/:value", h.Update)
	e.Any("/*", func(c echo.Context) error {
		return c.String(http.StatusBadRequest, "Bad request")
	})
	// start server
	e.Logger.Fatal(e.Start(serverAddress))
}
