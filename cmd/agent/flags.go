package main

import (
	"errors"
	"flag"
)

var serverAddress string
var reportInterval int64
var pollInterval int64

func parseFlags() error {
	flag.StringVar(&serverAddress, "a", "localhost:8080", "Server address to connect to")
	flag.Int64Var(&reportInterval, "r", 10, "Interval to report metrics")
	flag.Int64Var(&pollInterval, "p", 2, "Interval to poll metrics")
	flag.Parse()
	if len(flag.Args()) > 0 {
		return errors.New("entered unknown args")
	}
	return nil
}
