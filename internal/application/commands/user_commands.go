package commands

type RegisterUserCommand struct {
	Username string
	Password string
}

type ChangePasswordCommand struct {
	UserID      string
	OldPassword string
	NewPassword string
}
