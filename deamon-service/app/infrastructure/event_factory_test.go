package infrastructure

import (
	"github.com/stretchr/testify/mock"
	"github.com/yakob-abada/delfare/deamon-service/domain"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSecurityEventFactory_CreateEvent(t *testing.T) {
	mockLogger := new(MockLogger)
	mockLogger.On("Error", domain.LogContext{}, "Failed to encrypt message", mock.Anything).Return(nil)
	factory := NewSecurityEventFactory("this-is-a-32-byte-encryption-key", mockLogger)

	event := factory.CreateEvent()

	// Verify Criticality
	assert.True(t, event.Criticality >= 0 && event.Criticality <= 10, "Criticality should be between 0 and 10")

	// Verify Timestamp
	_, err := time.Parse(time.RFC3339, event.Timestamp)
	assert.NoError(t, err, "Timestamp should be in RFC3339 format")

	// Verify Message
	assert.NotContains(t, event.Message, "Security alert", "Event message should contain 'Security alert'")
}
