package infrastructure

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/yakob-abada/delfare/reader-service/domain"
)

// MockEventRepository mocks the EventRepository interface
type MockEventRepository struct {
	mock.Mock
}

func (m *MockEventRepository) GetCriticalEvents(ctx context.Context, limit int, minCriticality int) ([]domain.Event, error) {
	args := m.Called(ctx, limit, minCriticality)
	return args.Get(0).([]domain.Event), args.Error(1)
}

type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Debug(ctx domain.LogContext, msg string, fields ...interface{}) {
	m.Called(ctx, msg, fields)
}

func (m *MockLogger) Info(ctx domain.LogContext, msg string, fields ...interface{}) {
	m.Called(ctx, msg, fields)
}

func (m *MockLogger) Warn(ctx domain.LogContext, msg string, fields ...interface{}) {
	m.Called(ctx, msg, fields)
}

func (m *MockLogger) Error(ctx domain.LogContext, msg string, fields ...interface{}) {
	m.Called(ctx, msg, fields)
}

func (m *MockLogger) Fatal(ctx domain.LogContext, msg string, fields ...interface{}) {
	m.Called(ctx, msg, fields)
}

// MockEventPublisher mocks the EventPublisher interface
type MockEventPublisher struct {
	mock.Mock
}

func (m *MockEventPublisher) Publish(event domain.Event) error {
	args := m.Called(event)
	return args.Error(0)
}
