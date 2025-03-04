package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/yakob-abada/delfare/writer-service/application"
	"github.com/yakob-abada/delfare/writer-service/config"
	"github.com/yakob-abada/delfare/writer-service/domain"
	"github.com/yakob-abada/delfare/writer-service/infrastructure"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	cfg := config.LoadConfig(ctx)
	isProd := cfg.Env == "prod"

	// Initialize logger
	logger, err := infrastructure.NewZapLogger(isProd)
	if err != nil {
		panic("Failed to initialize logger")
	}
	logger.Info(domain.LogContext{}, "Writer Service is starting...")

	influxRepo := infrastructure.NewInfluxDBRepository(*cfg, logger)

	nc := infrastructure.NewNATSClient(cfg.NATSURL, cfg.NATSUsername, cfg.NATSPassword)
	defer nc.Close()

	subscriber := infrastructure.NewNATSSubscriber(nc, logger)

	eventService := application.NewEventService(subscriber, influxRepo, logger)
	err = eventService.HandleEvent(ctx)
	if err != nil {
		logger.Fatal(domain.LogContext{}, "Error handling events", "error", err)
	}

	// Keep the service running to listen to NATS
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
	cancel()

	// Wait for context cancellation before exiting
	<-ctx.Done()
}
