package shared

type AggregateRoot interface {
	GetID() string
	GetVersion() int
	GetUncommittedChanges() []Event
	ClearUncommittedChanges()
	Apply(Event)
}

type BaseAggregateRoot struct {
	ID      string
	Version int
	Changes []Event
}

func (a *BaseAggregateRoot) GetID() string                  { return a.ID }
func (a *BaseAggregateRoot) GetVersion() int                { return a.Version }
func (a *BaseAggregateRoot) GetUncommittedChanges() []Event { return a.Changes }
func (a *BaseAggregateRoot) ClearUncommittedChanges()       { a.Changes = []Event{} }
