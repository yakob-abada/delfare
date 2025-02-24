package domain

import "context"

// EventRepository defines business logic for querying security events
type EventRepository interface {
	GetCriticalEvents(ctx context.Context, limit int, minCriticality int) ([]Event, error)
}
