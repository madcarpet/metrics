package main

import (
	"errors"
	"flag"
	"os"
	"strconv"
)

var serverAddress string
var reportInterval int64
var pollInterval int64

func parseFlags() error {
	var errRprtInterval, errPollInterval error
	flag.StringVar(&serverAddress, "a", "localhost:8080", "Server address to connect to")
	flag.Int64Var(&reportInterval, "r", 10, "Interval to report metrics")
	flag.Int64Var(&pollInterval, "p", 2, "Interval to poll metrics")
	flag.Parse()
	if len(flag.Args()) > 0 {
		return errors.New("entered unknown args")
	}
	if envSrvAddr := os.Getenv("ADDRESS"); envSrvAddr != "" {
		serverAddress = envSrvAddr
	}
	if envRprtInterval := os.Getenv("REPORT_INTERVAL"); envRprtInterval != "" {
		reportInterval, errRprtInterval = strconv.ParseInt(envRprtInterval, 10, 64)
		if errRprtInterval != nil {
			return errors.New("bad REPORT_INTERVAL variable")
		}
	}
	if envPollInterval := os.Getenv("POLL_INTERVAL"); envPollInterval != "" {
		pollInterval, errPollInterval = strconv.ParseInt(envPollInterval, 10, 64)
		if errPollInterval != nil {
			return errors.New("bad POLL_INTERVAL variable")
		}
	}
	return nil
}
