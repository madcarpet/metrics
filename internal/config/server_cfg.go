package config

import (
	"context"
	"errors"
	"flag"
	"os"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/madcarpet/metrics/internal/adapter/storage"
	"github.com/madcarpet/metrics/internal/adapter/storage/filestorage"
	"github.com/madcarpet/metrics/internal/adapter/storage/memstorage"
	"github.com/madcarpet/metrics/internal/adapter/storage/pgstorage"
	"github.com/madcarpet/metrics/internal/service/metrics"
)

// ServerConfig struct to keep main application parameters.
type ServerConfig struct {
	ServerAddress string
	LoggingLevel  string
	StoreInterval int64
	SyncWrite     bool
	FilePath      string
	IsRestore     bool
	Storage       storage.Repository
	DBUrl         string
	Key           string
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

// NewServerConfig creates server configuration.
func NewServerConfig() (*ServerConfig, error) {
	config := ServerConfig{}
	// server address.
	flag.StringVar(&config.ServerAddress, "a", "localhost:8080", "Address server listen to")
	// server logging level.
	flag.StringVar(&config.LoggingLevel, "l", "info", "Logging level")
	// interval to store metrics.
	flag.Int64Var(&config.StoreInterval, "i", 300, "Store interval")
	// path to store server data.
	flag.StringVar(&config.FilePath, "f", "/tmp/metrics-db.json", "Path to store server data")
	// flag is nessesary to restore DB from file.
	flag.BoolVar(&config.IsRestore, "r", false, "Restore DB from file")
	// DB address in DSN format.
	flag.StringVar(&config.DBUrl, "d", "", "DB Url or params in DSN format")
	// secret key for signature.
	flag.StringVar(&config.Key, "k", "", "Key for data signature")
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
	switch {
	case config.DBUrl != "":
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
	case config.IsRestore || len(config.FilePath) > 0:
		var err error
		config.Storage, err = filestorage.NewFileStorage(config.FilePath, config.SyncWrite)
		if err != nil {
			return nil, err
		}
	default:
		config.Storage = memstorage.NewMemStorage()
	}

	if keyEnv := os.Getenv("KEY"); keyEnv != "" {
		config.Key = keyEnv
	}

	config.Root = metrics.NewGetAllMetricsSvc(config.Storage)
	config.Value = metrics.NewGetMetricSvc(config.Storage)
	config.Update = metrics.NewUpdateMetricSvc(config.Storage)
	config.Updates = metrics.NewUpdateMetricsSvc(config.Storage)
	config.Ping = metrics.NewPingSvc(config.Storage)
	config.Router = echo.New()
	return &config, nil
}
