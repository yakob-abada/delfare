package domain

import "context"

type EventRepository interface {
	GetLastCriticalEvents(ctx context.Context, doneCh chan struct{}, getEvent func(event Event), criticalityThreshold, limit int) error
}
