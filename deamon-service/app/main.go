package main

import (
	"fmt"
	"github.com/yakob-abada/delfare/deamon-service/infrastucture/validation"
	"os"
	"time"

	"github.com/yakob-abada/delfare/deamon-service/application"
	infrastructure "github.com/yakob-abada/delfare/deamon-service/infrastucture"
)

func main() {
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		panic("NATS_URL environment variable is not set")
	}

	natsUsername := os.Getenv("NATS_USERNAME")
	if natsURL == "" {
		panic("NATS_USERNAME environment variable is not set")
	}

	natsPassword := os.Getenv("NATS_PASSWORD")
	if natsURL == "" {
		panic("NATS_PASSWORD environment variable is not set")
	}

	encryptionKey := os.Getenv("ENCRYPTION_KEY")
	if encryptionKey == "" {
		panic("ENCRYPTION_KEY environment variable is not set")
	}

	nc := infrastructure.NewNATSClient(natsURL, natsUsername, natsPassword)
	defer nc.Close()

	eventService := application.NewEventService(
		infrastructure.NewNATSPublisher(nc),
		validation.NewSecurityEventValidator(),
		infrastructure.NewSecurityEventFactory(encryptionKey),
	)

	for {
		if err := eventService.PublishEvent(); err != nil {
			fmt.Println("Error publishing event:", err)
		}

		time.Sleep(2 * time.Second)
	}
}
