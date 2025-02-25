package application

import (
	"context"
	"errors"
	"fmt"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/yakob-abada/delfare/writer-service/domain"
	"github.com/yakob-abada/delfare/writer-service/infrastructure"
)

func TestEventService_GetLastCriticalEvents(t *testing.T) {
	ctx := context.TODO()
	event := domain.Event{Criticality: 8, Timestamp: time.Now(), Message: "Mock Event 1"}
	mockRepo := &infrastructure.MockEventRepository{}
	mockRepo.On("Write", ctx, event).Return(nil, errors.New("something went wrong"))

	mockLogger := &infrastructure.MockLogger{}
	mockLogger.On("Info", domain.LogContext{}, "Received message", mock.Anything).Return()

	mockSub := &infrastructure.MockSubscriber{}
	mockSub.On("ProcessEvents", ctx, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		// Retrieve the callback function
		callback := args.Get(1).(func(domain.Event) error)
		callback(event) // Simulate event reception

	})

	eventService := NewEventService(mockSub, mockRepo, mockLogger)

	err := eventService.HandleEvent(ctx)

	assert.NoError(t, err)
}

func TestEventService_GetLastCriticalEvents_Error(t *testing.T) {
	ctx := context.TODO()
	event := domain.Event{Criticality: 8, Timestamp: time.Now(), Message: "Mock Event 1"}
	mockRepo := &infrastructure.MockEventRepository{}

	mockRepo.On("Write", ctx, event).Return(nil, errors.New("something went wrong"))

	mockLogger := &infrastructure.MockLogger{}
	mockLogger.On("Info", domain.LogContext{}, "Received message", mock.Anything).Return()
	mockLogger.On("Error", domain.LogContext{}, "Failed to get critical events", mock.Anything).Return()

	mockSub := &infrastructure.MockSubscriber{}
	mockSub.On("ProcessEvents", ctx, mock.Anything).Return(fmt.Errorf("something went wrong")).Run(func(args mock.Arguments) {
		// Retrieve the callback function
		callback := args.Get(1).(func(domain.Event) error)
		callback(event) // Simulate event reception
	})

	eventService := NewEventService(mockSub, mockRepo, mockLogger)

	err := eventService.HandleEvent(ctx)

	assert.EqualError(t, err, "failed to get critical events: something went wrong")
}
