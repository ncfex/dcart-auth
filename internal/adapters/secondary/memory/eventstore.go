package memory

import (
	"context"
	"fmt"
	"sync"

	"github.com/ncfex/dcart-auth/internal/domain/shared"
	"github.com/ncfex/dcart-auth/internal/domain/user"
)

type InMemoryEventStore struct {
	events map[string][]shared.Event
	mu     sync.RWMutex
}

func NewInMemoryEventStore() *InMemoryEventStore {
	return &InMemoryEventStore{
		events: make(map[string][]shared.Event),
	}
}

func (s *InMemoryEventStore) SaveEvents(ctx context.Context, aggregateID string, events []shared.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	currentEvents := s.events[aggregateID]

	if len(currentEvents) > 0 {
		lastEvent := currentEvents[len(currentEvents)-1]
		firstNewEvent := events[0]
		if firstNewEvent.GetVersion() != lastEvent.GetVersion()+1 {
			return fmt.Errorf("concurrent modification detected for aggregate %s", aggregateID)
		}
	}

	if _, exists := s.events[aggregateID]; !exists {
		s.events[aggregateID] = []shared.Event{}
	}
	s.events[aggregateID] = append(s.events[aggregateID], events...)

	return nil
}

func (s *InMemoryEventStore) GetEvents(ctx context.Context, aggregateID string) ([]shared.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	events, exists := s.events[aggregateID]
	if !exists {
		return []shared.Event{}, nil
	}

	eventsCopy := make([]shared.Event, len(events))
	copy(eventsCopy, events)

	return eventsCopy, nil
}

func (s *InMemoryEventStore) GetEventsByType(ctx context.Context, aggregateType string) ([]shared.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []shared.Event

	for _, events := range s.events {
		for _, event := range events {
			if event.GetAggregateType() == aggregateType {
				result = append(result, event)
			}
		}
	}

	return result, nil
}

func (s *InMemoryEventStore) GetEventsByUsername(ctx context.Context, username string) ([]shared.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, events := range s.events {
		for _, event := range events {
			if event.GetEventType() == "USER_REGISTERED" {
				if payload, ok := event.GetPayload().(user.UserRegisteredEventPayload); ok {
					if payload.Username == username {
						return s.events[event.GetAggregateID()], nil
					}
				}
			}
		}
	}

	return []shared.Event{}, nil
}

// testing
func (s *InMemoryEventStore) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.events = make(map[string][]shared.Event)
}

func (s *InMemoryEventStore) GetAllEvents() []shared.Event {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var allEvents []shared.Event
	for _, events := range s.events {
		allEvents = append(allEvents, events...)
	}
	return allEvents
}
