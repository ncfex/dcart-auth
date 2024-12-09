package shared

type EventType string

type EventRegistry interface {
	CreateEvent(eventType EventType) (Event, bool)
	RegisterEvent(eventType EventType, factory func() Event)
}

type eventRegistry struct {
	factories map[EventType]func() Event
}

func NewEventRegistry() EventRegistry {
	return &eventRegistry{
		factories: make(map[EventType]func() Event),
	}
}

func (r *eventRegistry) CreateEvent(eventType EventType) (Event, bool) {
	factory, exists := r.factories[eventType]
	if !exists {
		return nil, false
	}
	return factory(), true
}

func (r *eventRegistry) RegisterEvent(eventType EventType, factory func() Event) {
	r.factories[eventType] = factory
}
