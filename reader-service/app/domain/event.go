package domain

import "time"

type Event struct {
	RequestID   string    `json:"request_id"`
	Criticality int       `json:"criticality"`
	Timestamp   time.Time `json:"timestamp"`
	Message     string    `json:"message"`
}

func (e Event) IsCritical(threshold int) bool {
	return e.Criticality > threshold
}
