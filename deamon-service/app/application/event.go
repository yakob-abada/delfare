package application

import (
	"fmt"

	"github.com/yakob-abada/delfare/deamon-service/domain"
	"github.com/yakob-abada/delfare/deamon-service/infrastucture"
	"github.com/yakob-abada/delfare/deamon-service/infrastucture/validation"
)

type EventService struct {
	publisher    domain.EventPublisher
	validator    validation.Validator
	eventFactory infrastructure.EventFactory
}

func NewEventService(
	Publisher domain.EventPublisher,
	validator validation.Validator,
	eventFactory infrastructure.EventFactory,
) *EventService {
	return &EventService{
		publisher:    Publisher,
		validator:    validator,
		eventFactory: eventFactory,
	}
}

func (s *EventService) PublishEvent() error {
	event := s.eventFactory.CreateEvent()

	// Validate the event
	if err := s.validator.Validate(event); err != nil {
		return fmt.Errorf("validation failed: %v", err)
	}

	// Publish the event
	if err := s.publisher.Publish(event); err != nil {
		return fmt.Errorf("failed to publish event: %v", err)
	}

	return nil
}
