package inbound

import (
	userDomain "github.com/ncfex/dcart-auth/internal/domain/user"
)

type UserSevice interface {
	CreateUser(username, password string) (*userDomain.User, error)
	ComparePassword(hashedPassword, password string) error
}
