package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yakob-abada/delfare/deamon-service/application"
	"github.com/yakob-abada/delfare/deamon-service/config"
	"github.com/yakob-abada/delfare/deamon-service/infrastructure"
	"github.com/yakob-abada/delfare/deamon-service/infrastructure/validation"
)

func main() {
	// Create a cancellable context that listens for OS signals
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := config.LoadConfig(ctx)
	isProd := cfg.Env == "prod"

	logger, err := infrastructure.NewZapLogger(isProd)
	if err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}

	encryptionKey := os.Getenv("ENCRYPTION_KEY")
	if encryptionKey == "" {
		log.Fatal("ENCRYPTION_KEY environment variable is not set")
	}

	nc := infrastructure.NewNATSClient(cfg)
	defer nc.Close()

	eventService := application.NewEventService(
		infrastructure.NewNATSPublisher(nc),
		validation.NewSecurityEventValidator(),
		infrastructure.NewSecurityEventFactory(encryptionKey, logger),
		logger,
	)

	// Start event publishing loop
	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Shutting down event publishing loop...")
				return
			default:
				if err := eventService.PublishEvent(); err != nil {
					fmt.Println("Error publishing event:", err)
				}
				time.Sleep(2 * time.Second)
			}
		}
	}()

	// Wait for termination signal
	<-ctx.Done()
	fmt.Println("Received shutdown signal, cleaning up...")

	// Allow some time for cleanup (if needed)
	time.Sleep(1 * time.Second)
	fmt.Println("Shutdown complete.")
}
