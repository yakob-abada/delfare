package application

import (
	"context"
	"fmt"

	"github.com/yakob-abada/delfare/writer-service/domain"
)

type EventService struct {
	subscriber domain.Subscriber
	repository domain.Repository
	logger     domain.Logger
}

func NewEventService(subscriber domain.Subscriber, repository domain.Repository, logger domain.Logger) *EventService {
	return &EventService{
		subscriber: subscriber,
		repository: repository,
		logger:     logger,
	}
}

func (s *EventService) HandleEvent(ctx context.Context) error {
	processor := func(event domain.Event) error {
		s.logger.Info(domain.LogContext{}, "Received message", "info", event)
		return s.repository.Write(ctx, event)
	}

	err := s.subscriber.ProcessEvents(ctx, processor)
	if err != nil {
		s.logger.Error(domain.LogContext{}, "Failed to get critical events", "error", err)
		return fmt.Errorf("failed to get critical events: %w", err)
	}

	return nil
}
