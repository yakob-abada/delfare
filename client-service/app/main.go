package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/yakob-abada/delfare/client-service/application"
	"github.com/yakob-abada/delfare/client-service/config"
	"github.com/yakob-abada/delfare/client-service/domain"
	"github.com/yakob-abada/delfare/client-service/infrastructure"
)

func main() {
	cfg := config.LoadConfig()
	isProd := cfg.Env == "prod"
	// Initialize Zap logger
	logger, err := infrastructure.NewZapLogger(isProd)
	if err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}

	defer logger.Close() // Flush logs before exiting

	logger.Info(domain.LogContext{}, "Client started")

	nc := infrastructure.NewNATSClient(cfg)
	defer nc.Close()

	repo := infrastructure.NewNATSEventRepository(nc, logger)
	eventService := application.NewEventService(repo, logger)
	ctx, cancel := context.WithCancel(context.Background())
	doneCh := make(chan struct{})

	// Fetch and display critical events
	err = eventService.GetLastCriticalEvents(ctx, doneCh, cfg.CriticalityThreshold, 5)
	if err != nil {
		logger.Fatal(domain.LogContext{}, "Error fetching critical events", "error", err)
	}

	// Handle graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		select {
		case <-sigChan:
			fmt.Println("Received termination signal. Exiting...")
			cancel()
		case <-doneCh:
			fmt.Println("Processing complete. Exiting...")
			cancel()
		}
	}()

	// Wait for context cancellation before exiting
	<-ctx.Done()
}
