package application

import (
	"fmt"

	"github.com/yakob-abada/delfare/deamon-service/domain"
	"github.com/yakob-abada/delfare/deamon-service/infrastructure"
	"github.com/yakob-abada/delfare/deamon-service/infrastructure/validation"
)

type EventService struct {
	publisher    domain.EventPublisher
	validator    validation.Validator
	eventFactory infrastructure.EventFactory
	logger       domain.Logger
}

func NewEventService(
	Publisher domain.EventPublisher,
	validator validation.Validator,
	eventFactory infrastructure.EventFactory,
	logger domain.Logger,
) *EventService {
	return &EventService{
		publisher:    Publisher,
		validator:    validator,
		eventFactory: eventFactory,
		logger:       logger,
	}
}

func (s *EventService) PublishEvent() error {
	event := s.eventFactory.CreateEvent()

	// Validate the event
	if err := s.validator.Validate(event); err != nil {
		s.logger.Error(domain.LogContext{RequestID: event.RequestID}, "validation failed", "error", err)
		return fmt.Errorf("validation failed: %v", err)
	}

	// Publish the event
	if err := s.publisher.Publish(event); err != nil {
		s.logger.Error(domain.LogContext{RequestID: event.RequestID}, "failed to publish event", "error", err)
		return fmt.Errorf("failed to publish event: %v", err)
	}

	return nil
}
