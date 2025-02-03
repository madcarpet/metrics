// Package config - app configuration.
package config

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/madcarpet/metrics/internal/adapter/storage"
	"github.com/madcarpet/metrics/internal/adapter/storage/filestorage"
	"github.com/madcarpet/metrics/internal/adapter/storage/memstorage"
	"github.com/madcarpet/metrics/internal/adapter/storage/pgstorage"
	"github.com/madcarpet/metrics/internal/parsers"
	"github.com/madcarpet/metrics/internal/service/metrics"
)

// ServerConfig struct to keep main application parameters.
type ServerConfig struct {
	ServerAddress       string `json:"address,omitempty"`
	LoggingLevel        string `json:"log_level,omitempty"`
	StoreInterval       int64
	SyncWrite           bool
	FilePath            string `json:"store_file,omitempty"`
	IsRestore           bool
	Storage             storage.Repository
	DBUrl               string `json:"database_dsn,omitempty"`
	Key                 string `json:"secret_key,omitempty"`
	PKeyPath            string `json:"crypto_key,omitempty"`
	CfgStoreInterval    string `json:"store_interval,omitempty"`
	CfgIsRestore        string `json:"restore,omitempty"`
	CfgStoreIntervalSet bool
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

var flagServerAddress string
var flagLoggingLevel string
var flagStoreInterval int64
var flagFilePath string
var flagIsRestore bool
var flagDBUrl string
var flagKey string
var flagPKeyPath string
var flagCfg string

// NewServerConfig creates server configuration.
func NewServerConfig() (*ServerConfig, error) {
	config := ServerConfig{}
	// server address.
	flag.StringVar(&flagServerAddress, "a", "", "Address server listen to")
	// server logging level.
	flag.StringVar(&flagLoggingLevel, "l", "", "Logging level")
	// interval to store metrics.
	flag.Int64Var(&flagStoreInterval, "i", 9999, "Store interval")
	// path to store server data.
	flag.StringVar(&flagFilePath, "f", "", "Path to store server data")
	// flag is nessesary to restore DB from file.
	flag.BoolVar(&flagIsRestore, "r", false, "Restore DB from file")
	// DB address in DSN format.
	flag.StringVar(&flagDBUrl, "d", "", "DB Url or params in DSN format")
	// secret key for signature.
	flag.StringVar(&flagKey, "k", "", "Key for data signature")
	// Path to private key.
	flag.StringVar(&flagPKeyPath, "crypto-key", "", "Path to secret key file for asymmetric encryption")
	// Configuration file name
	flag.StringVar(&flagCfg, "c", "", "Configuration file name")
	flag.Parse()

	if len(flag.Args()) > 0 {
		return nil, errors.New("entered unknown args")
	}

	// Get cfg data from file if c flag is set.
	if envCfg := os.Getenv("CONFIG"); envCfg != "" {
		flagCfg = envCfg
	}

	if flagCfg != "" {
		curPath, err := os.Getwd()
		if err != nil {
			return &config, errors.New("getting current directory problem")
		}
		cfgPath := filepath.Join(curPath, "cmd", "server", flagCfg)
		cfgData, err := os.ReadFile(cfgPath)
		if err != nil {
			fmt.Println("error")
			return &config, errors.New("cfg file read error")
		}
		err = json.Unmarshal(cfgData, &config)
		if err != nil {
			return &config, errors.New("cfg unmarshal error")
		}
		if config.CfgStoreInterval != "" {
			config.CfgStoreIntervalSet = true
			config.StoreInterval, err = parsers.ParseInterval(config.CfgStoreInterval)
			if err != nil {
				return &config, errors.New("wrong store interval or no store interval in config file")
			}
		}

		switch {
		case config.CfgIsRestore == "true":
			config.IsRestore = true
		default:
			config.IsRestore = false
		}
	}

	// Get environment variables.
	envSrvAddr := os.Getenv("ADDRESS")
	envLogLevel := os.Getenv("LOGLEVEL")
	envStoreInt := os.Getenv("STORE_INTERVAL")
	envFileStorePath := os.Getenv("FILE_STORAGE_PATH")
	envIsRstr := os.Getenv("RESTORE")
	envDBURL := os.Getenv("DATABASE_DSN")
	envKey := os.Getenv("KEY")
	envPKeyPath := os.Getenv("CRYPTO-KEY")

	// Set parameters according to the order
	switch {
	case envSrvAddr != "":
		config.ServerAddress = envSrvAddr
	case flagServerAddress != "":
		config.ServerAddress = flagServerAddress
	case config.ServerAddress != "":
	default:
		config.ServerAddress = "localhost:8080"
	}

	switch {
	case envLogLevel != "":
		config.LoggingLevel = envLogLevel
	case flagLoggingLevel != "":
		config.LoggingLevel = flagLoggingLevel
	case config.LoggingLevel != "":
	default:
		config.LoggingLevel = "info"
	}

	switch {
	case envStoreInt != "":
		i, err := strconv.ParseInt(envStoreInt, 10, 64)
		if err != nil {
			return nil, errors.New("store interval in envar STORE_INTERVAL is incorrect")
		}
		config.StoreInterval = i
	case flagStoreInterval != 9999:
		config.StoreInterval = flagStoreInterval
	case config.CfgStoreIntervalSet:
	default:
		config.StoreInterval = 300
	}

	switch config.StoreInterval {
	case 0:
		config.SyncWrite = true
	default:
		config.SyncWrite = false
	}

	switch {
	case envFileStorePath != "":
		config.FilePath = envFileStorePath
	case flagFilePath != "":
		config.FilePath = flagFilePath
	case config.FilePath != "":
	default:
		config.FilePath = "/tmp/metrics-db.json"
	}

	switch {
	case envIsRstr != "":
		isRstrValue, err := strconv.ParseBool(envIsRstr)
		if err != nil {
			return nil, errors.New("restore value in envar RESTORE is incorrect")
		}
		config.IsRestore = isRstrValue
	case flagIsRestore:
		config.IsRestore = true
	case config.IsRestore:
	default:
		config.IsRestore = false
	}

	switch {
	case envDBURL != "":
		config.DBUrl = envDBURL
	case flagDBUrl != "":
		config.DBUrl = flagDBUrl
	case config.DBUrl != "":
	default:
		config.DBUrl = ""
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

	switch {
	case envKey != "":
		config.Key = envKey
	case flagKey != "":
		config.Key = flagKey
	case config.Key != "":
	default:
		config.Key = ""
	}

	switch {
	case envPKeyPath != "":
		config.PKeyPath = envPKeyPath
	case flagPKeyPath != "":
		config.PKeyPath = flagPKeyPath
	case config.PKeyPath != "":
	default:
		config.PKeyPath = ""
	}

	config.Root = metrics.NewGetAllMetricsSvc(config.Storage)
	config.Value = metrics.NewGetMetricSvc(config.Storage)
	config.Update = metrics.NewUpdateMetricSvc(config.Storage)
	config.Updates = metrics.NewUpdateMetricsSvc(config.Storage)
	config.Ping = metrics.NewPingSvc(config.Storage)
	config.Router = echo.New()
	return &config, nil
}
