package user

import (
	"github.com/ncfex/dcart-auth/internal/domain/shared"
)

type UserFactory struct{}

func NewUserFactory() *UserFactory {
	return &UserFactory{}
}

func (f *UserFactory) CreateEmpty(id string) *User {
	return &User{
		BaseAggregateRoot: shared.BaseAggregateRoot{
			ID:      id,
			Version: 0,
			Changes: []shared.Event{},
		},
	}
}
