package config

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/madcarpet/metrics/internal/adapter/storage"
	"github.com/madcarpet/metrics/internal/adapter/storage/filestorage"
	"github.com/madcarpet/metrics/internal/adapter/storage/memstorage"
	"github.com/madcarpet/metrics/internal/adapter/storage/pgstorage"
	"github.com/madcarpet/metrics/internal/handlers/httpecho"
	"github.com/madcarpet/metrics/internal/service/metrics"
)

type ServerConfig struct {
	ServerAddress string
	LoggingLevel  string
	StoreInterval int64
	SyncWrite     bool
	FilePath      string
	IsRestore     bool
	Storage       storage.Repository
	DBUrl         string
	Services
	Router *echo.Echo
}

type Services struct {
	Root    *metrics.GetAllMetricsSvc
	Value   *metrics.GetMetricSvc
	Update  *metrics.UpdateMetricSvc
	Updates *metrics.UpdateMetricsSvc
	Ping    *metrics.PingSvc
}

func NewServerConfig() (*ServerConfig, error) {
	config := ServerConfig{}
	flag.StringVar(&config.ServerAddress, "a", "localhost:8080", "Address server listen to")
	flag.StringVar(&config.LoggingLevel, "l", "info", "Logging level")
	flag.Int64Var(&config.StoreInterval, "i", 300, "Store interval")
	flag.StringVar(&config.FilePath, "f", "/tmp/metrics-db.json", "Path to store server data")
	flag.BoolVar(&config.IsRestore, "r", false, "Restore DB from file")
	flag.StringVar(&config.DBUrl, "d", "", "DB Url or params in DSN format")
	flag.Parse()
	if len(flag.Args()) > 0 {
		return nil, errors.New("entered unknown args")
	}
	if envSrvAddr := os.Getenv("ADDRESS"); envSrvAddr != "" {
		config.ServerAddress = envSrvAddr
	}
	if loglvl := os.Getenv("LOGLEVEL"); loglvl != "" {
		config.LoggingLevel = loglvl
	}
	if storeInt := os.Getenv("STORE_INTERVAL"); storeInt != "" {
		i, err := strconv.ParseInt(storeInt, 10, 64)
		if err != nil {
			return nil, errors.New("store interval in envar STORE_INTERVAL is incorrect")
		}
		config.StoreInterval = i
	}
	switch config.StoreInterval {
	case 0:
		config.SyncWrite = true
	default:
		config.SyncWrite = false
	}

	if fileStorePath := os.Getenv("FILE_STORAGE_PATH"); fileStorePath != "" {
		config.FilePath = fileStorePath
	}
	if isRstr := os.Getenv("RESTORE"); isRstr != "" {
		isRstrValue, err := strconv.ParseBool(isRstr)
		if err != nil {
			return nil, errors.New("restore value in envar RESTORE is incorrect")
		}
		config.IsRestore = isRstrValue
	}

	if dbURL := os.Getenv("DATABASE_DSN"); dbURL != "" {
		config.DBUrl = dbURL
	}

	if config.DBUrl != "" {
		var err error
		pg, err := pgstorage.NewPGStorage(config.DBUrl)
		if err != nil {
			return nil, err
		}
		ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
		defer cancel()
		if err = pg.DB.PingContext(ctx); err == nil {
			err = pgstorage.DBMigration("./migrations", pg.DB)
			if err != nil {
				return nil, err
			}
		}
		config.Storage = pg
	} else if config.IsRestore || len(config.FilePath) > 0 {
		var err error
		config.Storage, err = filestorage.NewFileStorage(config.FilePath, config.SyncWrite)
		if err != nil {
			return nil, err
		}
	} else {
		config.Storage = memstorage.NewMemStorage()
	}
	config.Root = metrics.NewGetAllMetricsSvc(config.Storage)
	config.Value = metrics.NewGetMetricSvc(config.Storage)
	config.Update = metrics.NewUpdateMetricSvc(config.Storage)
	config.Updates = metrics.NewUpdateMetricsSvc(config.Storage)
	config.Ping = metrics.NewPingSvc(config.Storage)
	config.Router = echo.New()
	return &config, nil
}

func (sc *ServerConfig) Start() error {
	fmt.Println(sc.FilePath, sc.StoreInterval)
	httpecho.SetupRouter(sc.Router, sc.Root, sc.Value, sc.Update, sc.Updates, sc.Ping)
	if sc.IsRestore {
		err := sc.Storage.ImportFromFile()
		if err != nil {
			return err
		}
	}
	if sc.StoreInterval > 0 && sc.FilePath != "" {
		go func() {
			for {
				sc.Storage.ExportToFile()
				time.Sleep(time.Duration(sc.StoreInterval) * time.Second)
			}
		}()
	}
	go sc.Router.Start(sc.ServerAddress)
	return nil
}

func (sc *ServerConfig) Stop() error {
	err := sc.Storage.ExportToFile()
	if err != nil {
		return err
	}
	err = sc.Storage.Close()
	if err != nil {
		return err
	}
	return nil
}
