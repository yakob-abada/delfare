package validation

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/yakob-abada/delfare/deamon-service/domain"
)

func TestSecurityEventValidator_Validate(t *testing.T) {
	validator := NewSecurityEventValidator()

	tests := []struct {
		name        string
		event       domain.Event
		expectedErr string
	}{
		{
			name: "Valid Event",
			event: domain.Event{
				Criticality: 5,
				Timestamp:   time.Now().Format(time.RFC3339),
				Message:     "Valid event",
			},
			expectedErr: "",
		},
		{
			name: "Invalid Criticality (Negative)",
			event: domain.Event{
				Criticality: -1,
				Timestamp:   time.Now().Format(time.RFC3339),
				Message:     "Invalid criticality",
			},
			expectedErr: "invalid criticality value: must be between 0 and 10",
		},
		{
			name: "Invalid Criticality (Too High)",
			event: domain.Event{
				Criticality: 11,
				Timestamp:   time.Now().Format(time.RFC3339),
				Message:     "Invalid criticality",
			},
			expectedErr: "invalid criticality value: must be between 0 and 10",
		},
		{
			name: "Empty Event Message",
			event: domain.Event{
				Criticality: 5,
				Timestamp:   time.Now().Format(time.RFC3339),
				Message:     "",
			},
			expectedErr: "event message cannot be empty",
		},
		{
			name: "Invalid Timestamp Format",
			event: domain.Event{
				Criticality: 5,
				Timestamp:   "invalid-timestamp",
				Message:     "Valid event",
			},
			expectedErr: "", // Assuming timestamp format validation is not part of the validator
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Validate(tt.event)
			if tt.expectedErr == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.expectedErr)
			}
		})
	}
}
