package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/yakob-abada/delfare/reader-service/application"
	"github.com/yakob-abada/delfare/reader-service/config"
	"github.com/yakob-abada/delfare/reader-service/domain"
	"github.com/yakob-abada/delfare/reader-service/infrastructure"
)

func main() {
	cfg := config.LoadConfig()
	isProd := cfg.Env == "prod"

	logger, err := infrastructure.NewZapLogger(isProd)
	if err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}

	influxRepo := infrastructure.NewInfluxDBRepository(*cfg, logger)

	nc := infrastructure.NewNATSClient(cfg.NATSURL, cfg.NATSUsername, cfg.NATSPassword)
	defer nc.Close()

	publisher := infrastructure.NewNATSPublisher(nc, logger)
	eventService := application.NewEventService(influxRepo, publisher, logger, cfg.WorkerCount)

	ctx, cancel := context.WithCancel(context.Background())

	err = eventService.PublishCriticalEvents(ctx, cfg.MaxProcessedEvents, cfg.MinCriticalityEvents)
	if err != nil {
		log.Printf("failed to process event: %v\n", err)
		logger.Fatal(domain.LogContext{}, "Error process events", "error", err)
	}

	// Keep the service running to listen to NATS
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
	cancel()

	// Wait for context cancellation before exiting
	<-ctx.Done()
}
