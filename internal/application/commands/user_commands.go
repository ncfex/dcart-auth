package commands

type AuthenticateUserCommand struct {
	Username string
	Password string
}

type RegisterUserCommand struct {
	Username string
	Password string
}

type ChangePasswordCommand struct {
	UserID      string
	OldPassword string
	NewPassword string
}
