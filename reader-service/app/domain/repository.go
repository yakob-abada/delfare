package domain

import "context"

// Repository defines business logic for querying security events
type Repository interface {
	GetCriticalEvents(ctx context.Context, limit int, minCriticality int) ([]Event, error)
}
