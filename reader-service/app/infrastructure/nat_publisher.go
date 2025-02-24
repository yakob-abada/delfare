package infrastructure

import (
	"encoding/json"
	"github.com/nats-io/nats.go"
	"github.com/yakob-abada/delfare/reader-service/domain"
)

type NATSPublisher struct {
	nc     *nats.Conn
	logger domain.Logger
}

func NewNATSPublisher(nc *nats.Conn, logger domain.Logger) *NATSPublisher {
	return &NATSPublisher{nc: nc, logger: logger}
}

func (p *NATSPublisher) Publish(event domain.Event) error {
	data, _ := json.Marshal(event)
	if err := p.nc.Publish("events", data); err != nil {
		p.logger.Error(domain.LogContext{RequestID: event.RequestID}, "Failed to publish event", "error", err)

		return err
	}

	p.logger.Info(domain.LogContext{RequestID: event.RequestID}, "Published event", "info", string(data))

	return nil
}
