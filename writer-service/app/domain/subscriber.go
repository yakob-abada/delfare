package domain

import "context"

type Subscriber interface {
	ProcessEvents(ctx context.Context, processor func(event Event) error) error
}
