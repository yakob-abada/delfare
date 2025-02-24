package infrastructure

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/yakob-abada/delfare/client-service/domain"
)

type MockEventRepository struct {
	mock.Mock
}

func (m *MockEventRepository) GetLastCriticalEvents(
	ctx context.Context, doneCh chan struct{}, getEvent func(event domain.Event), threshold int, limit int,
) error {
	args := m.Called(ctx, doneCh, getEvent, threshold, limit)

	// Ensure correct return value count (error is expected as the second return value)
	if err, ok := args.Get(0).(error); ok {
		return err
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
