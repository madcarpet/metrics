package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/madcarpet/metrics/internal/config"
	"github.com/madcarpet/metrics/internal/logger"
)

func main() {
	//create channels for error and stopping
	sigChan := make(chan os.Signal, 1)
	errChan := make(chan error)
	//register system signals with channels
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	serverConfig, err := config.NewServerConfig()
	if err != nil {
		fmt.Printf("Server starting error: %v\n", err)
		os.Exit(1)
	}

	logger.Initialize(serverConfig.LoggingLevel)
	defer logger.Log.Sync()

	logger.Log.Info("Server starting")
	serverConfig.Start()

	//handle channels
	select {
	case stop := <-sigChan:
		fmt.Printf("Server stopping, recieved signal: %v\n", stop)
		serverConfig.Stop()
	case err := <-errChan:
		if err != nil {
			fmt.Printf("Server starting error: %v\n", err)
			os.Exit(1)
		}

	}
}
