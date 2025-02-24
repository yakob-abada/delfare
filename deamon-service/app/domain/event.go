package domain

type Event struct {
	RequestID   string `json:"request_id"`
	Criticality int    `json:"criticality"`
	Timestamp   string `json:"timestamp"`
	Message     string `json:"message"`
}
