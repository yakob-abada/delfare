package domain

import "context"

type Repository interface {
	Write(ctx context.Context, event Event) error
}
