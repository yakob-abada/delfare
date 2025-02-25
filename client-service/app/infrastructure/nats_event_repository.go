package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	"sync/atomic"

	"github.com/nats-io/nats.go"
	"github.com/yakob-abada/delfare/client-service/domain"
)

type NATSEventRepository struct {
	nc     *nats.Conn
	logger domain.Logger
}

func NewNATSEventRepository(nc *nats.Conn, logger domain.Logger) *NATSEventRepository {
	return &NATSEventRepository{nc: nc, logger: logger}
}

func (r *NATSEventRepository) GetLastCriticalEvents(
	ctx context.Context, doneCh chan struct{}, getEvent func(event domain.Event), criticalityThreshold, limit int,
) error {
	var count int32 = 0
	var sub *nats.Subscription

	sub, err := r.nc.Subscribe("events", func(msg *nats.Msg) {
		// Check if context is canceled before processing
		select {
		case <-ctx.Done():
			fmt.Println("Context canceled, stopping subscription.")
			return
		default:
		}

		var event domain.Event
		if err := json.Unmarshal(msg.Data, &event); err != nil {
			r.logger.Error(domain.LogContext{}, "Error unmarshalling event", "error", err)
			return
		}

		if event.Criticality >= criticalityThreshold {
			atomic.AddInt32(&count, 1)
			getEvent(event)
		}

		// Stop after limit matching events.
		if atomic.LoadInt32(&count) >= int32(limit) {
			r.logger.Info(domain.LogContext{}, fmt.Sprintf("Received %d critical events. Unsubscribing...", count))

			if err := sub.Unsubscribe(); err != nil {
				r.logger.Error(domain.LogContext{}, "Failed to unsubscribe:", err)
			}

			close(doneCh)
		}
	})

	return err
}
