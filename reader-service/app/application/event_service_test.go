package application

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/yakob-abada/delfare/reader-service/domain"
	"github.com/yakob-abada/delfare/reader-service/infrastructure"
)

func TestEventService_PublishCriticalEvents_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(infrastructure.MockEventRepository)
	mockPublisher := new(infrastructure.MockEventPublisher)
	mockLogger := new(infrastructure.MockLogger)
	events := []domain.Event{{Criticality: 8, Message: "Mock Event 1"}, {Criticality: 9, Message: "Mock Event 2"}}

	mockRepo.On("GetCriticalEvents", ctx, 2, 5).Return(events, nil)
	mockPublisher.On("Publish", mock.Anything).Return(nil)

	eventService := NewEventService(mockRepo, mockPublisher, mockLogger, 2)
	err := eventService.PublishCriticalEvents(ctx, 2, 5)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockPublisher.AssertExpectations(t)
}

func TestEventService_PublishCriticalEvents_FailureFetchingEvents(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(infrastructure.MockEventRepository)
	mockPublisher := new(infrastructure.MockEventPublisher)
	mockLogger := new(infrastructure.MockLogger)
	expectedErr := errors.New("failed to fetch events")

	mockRepo.On("GetCriticalEvents", ctx, 2, 5).Return([]domain.Event{}, expectedErr)
	mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything).Return()

	eventService := NewEventService(mockRepo, mockPublisher, mockLogger, 2)
	err := eventService.PublishCriticalEvents(ctx, 2, 5)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	mockRepo.AssertExpectations(t)
	mockLogger.AssertExpectations(t) // Ensure the error log is triggered
}

func TestEventService_PublishCriticalEvents_FailurePublishing(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(infrastructure.MockEventRepository)
	mockPublisher := new(infrastructure.MockEventPublisher)
	mockLogger := new(infrastructure.MockLogger)
	events := []domain.Event{{Criticality: 8, Message: "Mock Event 1"}}
	expectedErr := errors.New("failed to publish event")

	mockRepo.On("GetCriticalEvents", ctx, 1, 5).Return(events, nil)
	mockPublisher.On("Publish", mock.Anything).Return(expectedErr)
	mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything).Return()

	eventService := NewEventService(mockRepo, mockPublisher, mockLogger, 1)
	err := eventService.PublishCriticalEvents(ctx, 1, 5)

	assert.NoError(t, err) // The function itself should not fail if publishing fails
	mockRepo.AssertExpectations(t)
	mockPublisher.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}
