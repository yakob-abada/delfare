package application

import (
	"context"
	"fmt"

	"github.com/yakob-abada/delfare/client-service/domain"
)

type EventService struct {
	repo   domain.EventRepository
	logger domain.Logger
}

func NewEventService(repo domain.EventRepository, logger domain.Logger) *EventService {
	return &EventService{repo: repo, logger: logger}
}

func (s *EventService) GetLastCriticalEvents(ctx context.Context, doneCh chan struct{}, eventCh chan domain.Event, threshold int, limit int) error {
	getEvent := func(event domain.Event) {
		eventCh <- event
		s.logger.Info(domain.LogContext{}, "Received message", "info", event)
		return
	}

	err := s.repo.GetLastCriticalEvents(ctx, doneCh, getEvent, threshold, limit)
	if err != nil {
		s.logger.Error(domain.LogContext{}, "Failed to get critical events", "error", err)
		return fmt.Errorf("failed to get critical events: %w", err)
	}

	return nil
}
