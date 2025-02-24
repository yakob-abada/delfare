package application_test

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/mock"
	infrastructure "github.com/yakob-abada/delfare/deamon-service/infrastucture"
	"github.com/yakob-abada/delfare/deamon-service/infrastucture/validation"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/yakob-abada/delfare/deamon-service/application"
	"github.com/yakob-abada/delfare/deamon-service/domain"
)

func TestEventService_PublishEvent(t *testing.T) {
	// Create mocks
	mockPublisher := new(infrastructure.MockNATSPublisher)
	mockValidator := new(validation.MockValidator)
	mockEventFactory := new(infrastructure.MockEventFactory)

	// Create the EventService with the mocks
	eventService := application.NewEventService(mockPublisher, mockValidator, mockEventFactory)

	// Create a sample event
	event := domain.Event{
		Criticality: 5,
		Timestamp:   time.Now().Format(time.RFC3339),
		Message:     "Security alert",
	}

	// Set expectations on the mocks
	mockEventFactory.On("CreateEvent").Return(event)
	mockValidator.On("Validate", event).Return(nil)
	mockPublisher.On("Publish", event).Return(nil)

	// Call the PublishEvent method
	err := eventService.PublishEvent()

	// Verify that there were no errors
	assert.NoError(t, err)

	// Verify that the mocks were called as expected
	mockEventFactory.AssertCalled(t, "CreateEvent")
	mockValidator.AssertCalled(t, "Validate", event)
	mockPublisher.AssertCalled(t, "Publish", event)
}

func TestEventService_PublishEvent_ValidationError(t *testing.T) {
	// Create mocks
	mockPublisher := new(infrastructure.MockNATSPublisher)
	mockValidator := new(validation.MockValidator)
	mockEventFactory := new(infrastructure.MockEventFactory)

	// Create the EventService with the mocks
	eventService := application.NewEventService(mockPublisher, mockValidator, mockEventFactory)

	// Create a sample event
	event := domain.Event{
		Criticality: 5,
		Timestamp:   time.Now().Format(time.RFC3339),
		Message:     "Security alert",
	}

	// Set expectations on the mocks
	mockEventFactory.On("CreateEvent").Return(event)
	mockValidator.On("Validate", event).Return(errors.New("validation failed"))

	// Call the PublishEvent method
	err := eventService.PublishEvent()

	// Verify that the error is returned
	assert.EqualError(t, err, "validation failed: validation failed")

	// Verify that the mocks were called as expected
	mockEventFactory.AssertCalled(t, "CreateEvent")
	mockValidator.AssertCalled(t, "Validate", event)
	mockPublisher.AssertNotCalled(t, "Publish", event)
}

func TestEventService_PublishEvent_Error(t *testing.T) {
	// Create a mock EventPublisher
	mockPublisher := new(infrastructure.MockNATSPublisher)
	mockValidator := new(validation.MockValidator)
	mockEventFactory := new(infrastructure.MockEventFactory)

	// Create a sample event
	event := domain.Event{
		Criticality: 5,
		Timestamp:   time.Now().Format(time.RFC3339),
		Message:     "Security alert",
	}

	// Set expectations on the mocks
	mockEventFactory.On("CreateEvent").Return(event)
	mockValidator.On("Validate", event).Return(nil)
	mockPublisher.On("Publish", mock.AnythingOfType("domain.Event")).Return(fmt.Errorf("publish error"))

	// Create the EventService with the mock publisher
	eventService := application.NewEventService(mockPublisher, mockValidator, mockEventFactory)

	// Call the PublishEvent method
	err := eventService.PublishEvent()
	assert.EqualError(t, err, "failed to publish event: publish error")

	// Verify that the Publish method was called
	mockPublisher.AssertCalled(t, "Publish", mock.AnythingOfType("domain.Event"))
}

func TestEventService_PublishEvent_MultipleCalls(t *testing.T) {
	// Create mocks
	mockPublisher := new(infrastructure.MockNATSPublisher)
	mockValidator := new(validation.MockValidator)
	mockEventFactory := new(infrastructure.MockEventFactory)

	// Create the EventService with the mocks
	eventService := application.NewEventService(mockPublisher, mockValidator, mockEventFactory)

	// Create a sample event
	event := domain.Event{
		Criticality: 5,
		Timestamp:   time.Now().Format(time.RFC3339),
		Message:     "Security alert",
	}

	// Set expectations on the mocks
	mockEventFactory.On("CreateEvent").Return(event)
	mockValidator.On("Validate", event).Return(nil)
	mockPublisher.On("Publish", mock.AnythingOfType("domain.Event")).Return()

	// Call the PublishEvent method multiple times
	for i := 0; i < 5; i++ {
		err := eventService.PublishEvent()

		// Verify that there were no errors
		assert.NoError(t, err)
	}

	// Verify that the Publish method was called 5 times
	mockPublisher.AssertNumberOfCalls(t, "Publish", 5)
}
