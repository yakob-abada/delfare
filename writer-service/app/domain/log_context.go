package domain

type LogContext struct {
	RequestID     string // Unique ID for the request
	CorrelationID string // ID to correlate logs across services
	UserID        string // Optional: ID of the user making the request
}
