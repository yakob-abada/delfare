package application

import (
	"context"
	"sync"

	"github.com/yakob-abada/delfare/reader-service/domain"
)

// EventService implements EventService
type EventService struct {
	repo             domain.Repository
	publisher        domain.Publisher
	logger           domain.Logger
	eventWorkerCount int
	mu               sync.RWMutex
	criticalityLevel int // Tracks current criticality level
}

func NewEventService(
	repo domain.Repository, publisher domain.Publisher, logger domain.Logger, eventWorkerCount int,
) *EventService {
	return &EventService{
		repo:             repo,
		publisher:        publisher,
		logger:           logger,
		eventWorkerCount: eventWorkerCount,
	}
}

// GetCriticality returns the latest recorded criticality level
func (s *EventService) GetCriticality() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.criticalityLevel
}

// PublishCriticalEvents fetches and publishes critical events
func (s *EventService) PublishCriticalEvents(ctx context.Context, limit int, minCriticality int) error {
	events, err := s.repo.GetCriticalEvents(ctx, limit, minCriticality)
	if err != nil {
		s.logger.Error(domain.LogContext{}, "Error fetching critical events", err)
		return err
	}

	// Update criticality level
	s.mu.Lock()
	s.criticalityLevel = len(events) // Assuming criticality is based on event count
	s.mu.Unlock()

	eventChan := make(chan domain.Event, len(events)) // Buffered channel
	var wg sync.WaitGroup

	// Start worker goroutines
	for i := 0; i < s.eventWorkerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for event := range eventChan {
				if err := s.publisher.Publish(event); err != nil {
					s.logger.Error(domain.LogContext{}, "Failed to publish event", err)
				}
			}
		}()
	}

	// Send events to the channel
	for _, event := range events {
		eventChan <- event
	}
	close(eventChan)

	// Wait for all workers to finish
	wg.Wait()

	return nil
}
