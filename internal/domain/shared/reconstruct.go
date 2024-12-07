package shared

import (
	"fmt"
	"sort"
)

type AggregateFactory[T AggregateRoot] interface {
	CreateEmpty(id string) T
}

type EventSourcedAggregate struct{}

func ReconstructAggregate[T AggregateRoot](
	events []Event,
	factory AggregateFactory[T],
) (T, error) {
	var empty T
	if len(events) == 0 {
		return empty, fmt.Errorf("no events found")
	}

	sortedEvents := make([]Event, len(events))
	copy(sortedEvents, events)
	sort.Slice(sortedEvents, func(i, j int) bool {
		return sortedEvents[i].GetVersion() < sortedEvents[j].GetVersion()
	})

	aggregate := factory.CreateEmpty(sortedEvents[0].GetAggregateID())

	for _, event := range sortedEvents {
		expectedVersion := aggregate.GetVersion() + 1
		if event.GetVersion() != expectedVersion {
			return empty, fmt.Errorf(
				"wrong event version: expected %d, got %d",
				expectedVersion,
				event.GetVersion(),
			)
		}

		aggregate.Apply(event)
	}

	return aggregate, nil
}
