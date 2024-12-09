// Package app - app start stop controller.
package app

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/madcarpet/metrics/internal/config"
	"github.com/madcarpet/metrics/internal/handlers/httpecho"
	"github.com/madcarpet/metrics/internal/logger"
	"go.uber.org/zap"
)

// App struct for app functionality, contains configuration.
type App struct {
	cfg *config.ServerConfig
}

// NewApp creates a new App.
func NewApp(c *config.ServerConfig) *App {
	return &App{
		cfg: c,
	}
}

// AppStart starts the app according to the configuration parameters.
func (a *App) AppStart(ctx context.Context) error {
	fmt.Println(a.cfg.FilePath, a.cfg.StoreInterval)
	if a.cfg.Key != "" {
		err := os.Setenv("SECRET_KEY", a.cfg.Key)
		if err != nil {
			return err
		}
		httpecho.SetupRouter(a.cfg.Router, a.cfg.Root, a.cfg.Value, a.cfg.Update, a.cfg.Updates, a.cfg.Ping, true)
	} else {
		httpecho.SetupRouter(a.cfg.Router, a.cfg.Root, a.cfg.Value, a.cfg.Update, a.cfg.Updates, a.cfg.Ping, false)
	}
	if a.cfg.IsRestore {
		err := a.cfg.Storage.ImportFromFile(ctx)
		if err != nil {
			return err
		}
	}
	if a.cfg.StoreInterval > 0 && a.cfg.FilePath != "" {
		go func() {
			for {
				err := a.cfg.Storage.ExportToFile(ctx)
				if err != nil {
					logger.Log.Error("error exporting to file", zap.Error(err))
					return
				}
				time.Sleep(time.Duration(a.cfg.StoreInterval) * time.Second)
			}
		}()
	}
	go func() {
		err := a.cfg.Router.Start(a.cfg.ServerAddress)
		if err != nil {
			logger.Log.Error("router start error", zap.Error(err))

		}
	}()
	return nil
}

// AppStop stops the app.
func (a *App) AppStop(ctx context.Context) error {
	err := a.cfg.Storage.ExportToFile(ctx)
	if err != nil {
		return err
	}
	err = a.cfg.Storage.Close()
	if err != nil {
		return err
	}
	return nil
}
