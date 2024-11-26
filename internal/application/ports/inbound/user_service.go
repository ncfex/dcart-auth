package inbound

import (
	"context"

	userDomain "github.com/ncfex/dcart-auth/internal/domain/user"
)

type UserSevice interface {
	CreateUser(ctx context.Context, username, password string) (*userDomain.User, error)
	ValidateWithCreds(ctx context.Context, username, password string) (*userDomain.User, error)
	ValidateWithID(ctx context.Context, userID string) (*userDomain.User, error)
}
