package app

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/madcarpet/metrics/internal/config"
	"github.com/madcarpet/metrics/internal/handlers/httpecho"
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
		go func() {
			for {
				a.cfg.Storage.ExportToFile(ctx)
				time.Sleep(time.Duration(a.cfg.StoreInterval) * time.Second)
			}
		}()
	}
	go a.cfg.Router.Start(a.cfg.ServerAddress)
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
