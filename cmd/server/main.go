package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/madcarpet/metrics/internal/app"
	"github.com/madcarpet/metrics/internal/config"
	"github.com/madcarpet/metrics/internal/logger"
)

func main() {
	// create channels for error and stopping.
	sigChan := make(chan os.Signal, 1)
	errChan := make(chan error)
	// register system signals with channels.
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	serverConfig, err := config.NewServerConfig()
	if err != nil {
		fmt.Printf("Server config preparing error: %v\n", err)
		os.Exit(1)
	}

	logger.Initialize(serverConfig.LoggingLevel)
	defer logger.Log.Sync()

	logger.Log.Info("Server starting")

	aplication := app.NewApp(serverConfig)
	err = aplication.AppStart(context.Background())
	if err != nil {
		fmt.Printf("Server starting error: %v\n", err)
		return
	}
	// handle channels.
	select {
	case stop := <-sigChan:
		fmt.Printf("Server stopping, recieved signal: %v\n", stop)
		err := aplication.AppStop(context.Background())
		if err != nil {
			fmt.Printf("server stopped wuth err %v", err)
		}
	case err := <-errChan:
		if err != nil {
			fmt.Printf("Server got error: %v\n", err)
			err := aplication.AppStop(context.Background())
			if err != nil {
				fmt.Printf("server stopped wuth err %v", err)
			}
		}

	}
}
