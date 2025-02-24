package application

import (
	"context"
	"errors"
	"testing"
	
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/yakob-abada/delfare/client-service/domain"
	"github.com/yakob-abada/delfare/client-service/infrastructure"
)

func TestGetLastCriticalEvents_Success(t *testing.T) {
	mockRepo := new(infrastructure.MockEventRepository)
	mockLogger := new(infrastructure.MockLogger)
	eventService := NewEventService(mockRepo, mockLogger)

	ctx := context.Background()
	doneCh := make(chan struct{})
	threshold := 5
	limit := 10
	testEvent := domain.Event{Criticality: 6, Message: "Critical Event"}

	mockRepo.On("GetLastCriticalEvents", ctx, doneCh, mock.Anything, threshold, limit).
		Return(nil).Run(func(args mock.Arguments) {
		// Retrieve the callback function
		callback := args.Get(2).(func(domain.Event))
		callback(testEvent) // Simulate event reception
	})

	mockLogger.On("Info", domain.LogContext{}, mock.Anything, mock.Anything).Return()

	err := eventService.GetLastCriticalEvents(ctx, doneCh, threshold, limit)
	require.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestGetLastCriticalEvents_Failure(t *testing.T) {
	mockRepo := new(infrastructure.MockEventRepository)
	mockLogger := new(infrastructure.MockLogger)
	eventService := NewEventService(mockRepo, mockLogger)

	ctx := context.Background()
	doneCh := make(chan struct{})
	threshold := 5
	limit := 10

	expectedError := errors.New("repository failure")
	mockRepo.On("GetLastCriticalEvents", ctx, doneCh, mock.Anything, threshold, limit).Return(expectedError)
	mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything).Return()

	err := eventService.GetLastCriticalEvents(ctx, doneCh, threshold, limit)
	require.Error(t, err)
	require.Equal(t, expectedError, errors.Unwrap(err))
	mockRepo.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}
