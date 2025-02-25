package infrastructure

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/yakob-abada/delfare/writer-service/domain"
)

type MockEventRepository struct {
	mock.Mock
}

func (m *MockEventRepository) Write(ctx context.Context, event domain.Event) error {
	args := m.Called(ctx, event)
	if args.Error(0) != nil {
		return args.Error(0)
	}

	return nil
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

type MockSubscriber struct {
	mock.Mock
}

func (m *MockSubscriber) ProcessEvents(ctx context.Context, processor func(event domain.Event) error) error {
	args := m.Called(ctx, processor)

	if err, ok := args.Get(0).(error); ok {
		return err
	}

	return nil
}
