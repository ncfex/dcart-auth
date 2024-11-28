package inbound

type RegisterRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UserDTO struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type TokenRequest struct {
	Token string `json:"token" validate:"required"`
}

type TokenPairDTO struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type TokenDTO struct {
	Token string `json:"token"`
}

type ValidateResponse struct {
	Valid bool    `json:"valid"`
	User  UserDTO `json:"user"`
}
