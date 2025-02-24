package infrastructure

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSecurityEventFactory_CreateEvent(t *testing.T) {
	factory := NewSecurityEventFactory()

	event := factory.CreateEvent()

	// Verify Criticality
	assert.True(t, event.Criticality >= 0 && event.Criticality <= 10, "Criticality should be between 0 and 10")

	// Verify Timestamp
	_, err := time.Parse(time.RFC3339, event.Timestamp)
	assert.NoError(t, err, "Timestamp should be in RFC3339 format")

	// Verify Message
	assert.Contains(t, event.Message, "Security alert", "Event message should contain 'Security alert'")
}
