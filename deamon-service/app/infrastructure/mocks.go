package infrastructure

import (
	"github.com/stretchr/testify/mock"
	"github.com/yakob-abada/delfare/deamon-service/domain"
)

type MockNATSPublisher struct {
	mock.Mock
}

func (m *MockNATSPublisher) Publish(event domain.Event) error {
	args := m.Called(event)
	return args.Error(0)
}

type MockEventFactory struct {
	mock.Mock
}

func (m *MockEventFactory) CreateEvent() domain.Event {
	args := m.Called()
	return args.Get(0).(domain.Event)
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
