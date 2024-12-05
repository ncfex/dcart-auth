package user

import (
	"errors"
	"fmt"
	"sort"

	"github.com/ncfex/dcart-auth/internal/domain/shared"
)

var (
	ErrCannotReconstructUser = errors.New("cannot reconstruct user")
)

func ReconstructFromEvents(events []shared.Event) (*User, error) {
	if len(events) == 0 {
		return nil, fmt.Errorf("%s: %w", ErrCannotReconstructUser.Error(), ErrUserNotFound)
	}

	sortedEvents := make([]shared.Event, len(events))
	copy(sortedEvents, events)
	sort.Slice(sortedEvents, func(i, j int) bool {
		return sortedEvents[i].GetVersion() < sortedEvents[j].GetVersion()
	})

	user := &User{
		BaseAggregateRoot: shared.BaseAggregateRoot{
			ID:      sortedEvents[0].GetAggregateID(),
			Version: 0,
			Changes: []shared.Event{},
		},
	}

	for _, event := range sortedEvents {
		user.Apply(event)
		user.Version = event.GetVersion()
	}

	return user, nil
}
