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

type App struct {
	cfg *config.ServerConfig
}

func NewApp(c *config.ServerConfig) *App {
	return &App{
		cfg: c,
	}
}

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
		go func() error {
			for {
				err := a.cfg.Storage.ExportToFile(ctx)
				if err != nil {
					logger.Log.Error("error exporting to file", zap.Error(err))
					return err
				}
				time.Sleep(time.Duration(a.cfg.StoreInterval) * time.Second)
			}
		}()
	}
	go func() error {
		err := a.cfg.Router.Start(a.cfg.ServerAddress)
		if err != nil {
			logger.Log.Error("router start error", zap.Error(err))
			return err
		}
		return nil
	}()
	return nil
}

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
