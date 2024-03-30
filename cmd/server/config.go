package main

import (
	"errors"
	"flag"
	"os"
	"strconv"
)

var serverAddress string
var loggingLevel string
var storeInterval int64
var fileStoragePath string
var isRestore bool

// function to parse args from cli or environment vars
func parseFlags() error {
	flag.StringVar(&serverAddress, "a", "localhost:8080", "Address server listen to")
	flag.StringVar(&loggingLevel, "l", "info", "Logging level")
	flag.Int64Var(&storeInterval, "i", 300, "Store interval")
	flag.StringVar(&fileStoragePath, "f", "/tmp/metrics-db.json", "Path to store server DB")
	flag.BoolVar(&isRestore, "r", false, "Restore DB from file")
	flag.Parse()
	if len(flag.Args()) > 0 {
		return errors.New("entered unknown args")
	}
	if envSrvAddr := os.Getenv("ADDRESS"); envSrvAddr != "" {
		serverAddress = envSrvAddr
	}
	if loglvl := os.Getenv("LOGLEVEL"); loglvl != "" {
		loggingLevel = loglvl
	}
	if storeInt := os.Getenv("STORE_INTERVAL"); storeInt != "" {
		i, err := strconv.ParseInt(storeInt, 10, 64)
		if err != nil {
			return errors.New("store interval in envar STORE_INTERVAL is incorrect")
		}
		storeInterval = i
	}
	if fileStorePath := os.Getenv("FILE_STORAGE_PATH"); fileStorePath != "" {
		fileStoragePath = fileStorePath
	}
	if isRstr := os.Getenv("RESTORE"); isRstr != "" {
		isRstrValue, err := strconv.ParseBool(isRstr)
		if err != nil {
			return errors.New("restore value in envar RESTORE is incorrect")
		}
		isRestore = isRstrValue
	}
	return nil
}
