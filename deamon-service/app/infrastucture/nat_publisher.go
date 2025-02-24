package infrastructure

import (
	"encoding/json"
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/yakob-abada/delfare/deamon-service/domain"
)

type NATSPublisher struct {
	nc *nats.Conn
}

func NewNATSPublisher(nc *nats.Conn) *NATSPublisher {
	return &NATSPublisher{nc: nc}
}

func (p *NATSPublisher) Publish(event domain.Event) error {
	data, _ := json.Marshal(event)
	if err := p.nc.Publish("events", data); err != nil {
		fmt.Println("Failed to publish event:", err)
	} else {
		fmt.Println("Published event:", string(data))
	}
	return nil
}
