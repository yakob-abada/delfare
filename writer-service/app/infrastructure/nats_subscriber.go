package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/yakob-abada/delfare/writer-service/domain"
)

type NATSSubscriber struct {
	nc     *nats.Conn
	logger domain.Logger
}

func NewNATSSubscriber(nc *nats.Conn, logger domain.Logger) *NATSSubscriber {
	return &NATSSubscriber{nc: nc, logger: logger}
}

func (s *NATSSubscriber) ProcessEvents(ctx context.Context, processor func(event domain.Event) error) error {
	_, err := s.nc.Subscribe("events", func(msg *nats.Msg) {
		// Check if context is canceled before processing
		select {
		case <-ctx.Done():
			fmt.Println("Context canceled, stopping subscription.")
			return
		default:
		}

		var event domain.Event
		if err := json.Unmarshal(msg.Data, &event); err != nil {
			ctx := domain.LogContext{RequestID: event.RequestID}
			s.logger.Error(ctx, "Error unmarshalling event", "error", err)
			return
		}

		err := processor(event)
		if err != nil {
			ctx := domain.LogContext{RequestID: event.RequestID}
			s.logger.Error(ctx, "failed to process", "error", err)
		}
	})

	//defer sub.Unsubscribe()

	return err
}
