package main

import (
	"errors"
	"flag"
)

var serverAddress string

func parseFlags() error {
	flag.StringVar(&serverAddress, "a", "localhost:8080", "Address server listen to")
	flag.Parse()
	if len(flag.Args()) > 0 {
		return errors.New("entered unknown args")
	}
	return nil
}
