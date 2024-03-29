package main

import (
	"errors"
	"flag"
	"os"
)

var serverAddress string

// function to parse args from cli or environment vars
func parseFlags() error {
	flag.StringVar(&serverAddress, "a", "localhost:8080", "Address server listen to")
	flag.Parse()
	if len(flag.Args()) > 0 {
		return errors.New("entered unknown args")
	}
	if envSrvAddr := os.Getenv("ADDRESS"); envSrvAddr != "" {
		serverAddress = envSrvAddr
	}
	return nil
}
