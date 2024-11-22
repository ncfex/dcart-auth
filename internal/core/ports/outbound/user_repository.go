package outbound

import (
	"context"

	"github.com/google/uuid"
	userDomain "github.com/ncfex/dcart-auth/internal/core/domain/user"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *userDomain.User) (*userDomain.User, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*userDomain.User, error)
	GetUserByUsername(ctx context.Context, username string) (*userDomain.User, error)
}
