package main

import (
	"encoding/json"
	"errors"
	"flag"
	"os"
	"path/filepath"
	"strconv"

	"github.com/madcarpet/metrics/internal/parsers"
)

type agentConfig struct {
	ServerAddress     string `json:"address,omitempty"`
	ReportInterval    int64
	PollInterval      int64
	SecretKey         string `json:"secret_key,omitempty"`
	RateLimit         int64  `json:"rate_limit,omitempty"`
	PKey              string `json:"crypto_key,omitempty"`
	CfgReportInterval string `json:"report_interval,omitempty"`
	CfgPollInterval   string `json:"poll_interval,omitempty"`
}

var flagServerAddress string
var flagReportInterval int64
var flagPollInterval int64
var flagSecretKey string
var flagRateLimit int64
var flagPKey string
var flagCfg string

// parseFlags function parses command line flags and environment variables for configuration settings.
func parseFlags() (agentConfig, error) {
	agentCfg := agentConfig{}
	var errRprtInterval, errPollInterval, errRateLimit error
	flag.StringVar(&flagServerAddress, "a", "", "Server address to connect to")
	flag.Int64Var(&flagReportInterval, "r", 0, "Interval to report metrics")
	flag.Int64Var(&flagPollInterval, "p", 0, "Interval to poll metrics")
	flag.StringVar(&flagSecretKey, "k", "", "Key for data signature")
	flag.Int64Var(&flagRateLimit, "l", 0, "Concurrent send rate limit")
	flag.StringVar(&flagPKey, "crypto-key", "", "Path to pub key for asymmetric encryption")
	flag.StringVar(&flagCfg, "c", "", "Config file name")
	flag.Parse()

	if len(flag.Args()) > 0 {
		return agentCfg, errors.New("entered unknown args")
	}

	// Get cfg data from file if c flag is set.
	if envCfg := os.Getenv("CONFIG"); envCfg != "" {
		flagCfg = envCfg
	}

	if flagCfg != "" {
		curPath, err := os.Getwd()
		if err != nil {
			return agentCfg, errors.New("getting current directory problem")
		}
		cfgPath := filepath.Join(curPath, flagCfg)
		cfgData, err := os.ReadFile(cfgPath)
		if err != nil {
			return agentCfg, errors.New("cfg file read error")
		}
		err = json.Unmarshal(cfgData, &agentCfg)
		if err != nil {
			return agentCfg, errors.New("cfg unmarshal error")
		}
		if agentCfg.CfgReportInterval != "" {
			agentCfg.ReportInterval, err = parsers.ParseInterval(agentCfg.CfgReportInterval)
			if err != nil {
				return agentCfg, errors.New("wrong report interval or no report interval in config file")
			}
		}
		if agentCfg.CfgPollInterval != "" {
			agentCfg.PollInterval, err = parsers.ParseInterval(agentCfg.CfgPollInterval)
			if err != nil {
				return agentCfg, errors.New("wrong poll interval or no poll interval in config file")
			}
		}
	}

	// Get environment variables.
	envSrvAddr := os.Getenv("ADDRESS")
	envRprtInterval := os.Getenv("REPORT_INTERVAL")
	envPollInterval := os.Getenv("POLL_INTERVAL")
	envSecretKey := os.Getenv("KEY")
	envRateLimit := os.Getenv("RATE_LIMIT")
	envPKey := os.Getenv("CRYPTO_KEY")

	// Set parameters according to the order
	switch {
	case envSrvAddr != "":
		agentCfg.ServerAddress = envSrvAddr
	case flagServerAddress != "":
		agentCfg.ServerAddress = flagServerAddress
	case agentCfg.ServerAddress != "":
	default:
		agentCfg.ServerAddress = "localhost:8080"
	}

	switch {
	case envRprtInterval != "":
		agentCfg.ReportInterval, errRprtInterval = strconv.ParseInt(envRprtInterval, 10, 64)
		if errRprtInterval != nil {
			return agentCfg, errors.New("bad REPORT_INTERVAL parameter")
		}
	case flagReportInterval != 0:
		agentCfg.ReportInterval = flagReportInterval
	case agentCfg.ReportInterval != 0:
	default:
		agentCfg.ReportInterval = 10
	}

	switch {
	case envPollInterval != "":
		agentCfg.PollInterval, errPollInterval = strconv.ParseInt(envPollInterval, 10, 64)
		if errPollInterval != nil {
			return agentCfg, errors.New("bad POLL_INTERVAL parameter")
		}
	case flagPollInterval != 0:
		agentCfg.PollInterval = flagPollInterval
	case agentCfg.PollInterval != 0:
	default:
		agentCfg.PollInterval = 2
	}

	switch {
	case envSecretKey != "":
		agentCfg.SecretKey = envSecretKey
	case flagSecretKey != "":
		agentCfg.SecretKey = flagSecretKey
	case agentCfg.SecretKey != "":
	default:
		agentCfg.SecretKey = ""
	}

	switch {
	case envRateLimit != "":
		agentCfg.RateLimit, errRateLimit = strconv.ParseInt(envRateLimit, 10, 64)
		if errRateLimit != nil {
			return agentCfg, errors.New("bad RATE_LIMIT parameter")
		}
	case flagRateLimit != 0:
		agentCfg.RateLimit = flagRateLimit
	case agentCfg.RateLimit != 0:
	default:
		agentCfg.RateLimit = 1
	}

	switch {
	case envPKey != "":
		agentCfg.PKey = envPKey
	case flagPKey != "":
		agentCfg.PKey = flagPKey
	case agentCfg.PKey != "":
	default:
		agentCfg.PKey = ""
	}

	return agentCfg, nil
}
