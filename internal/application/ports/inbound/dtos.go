package inbound

type RegisterRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type RefreshRequest struct {
	TokenString string `json:"token_string" validate:"required"`
}

type LogoutRequest struct {
	TokenString string `json:"token_string" validate:"required"`
}

type ValidateRequest struct {
	TokenString string `json:"token_string" validate:"required"`
}

type UserResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type ValidateResponse struct {
	Valid bool         `json:"valid"`
	User  UserResponse `json:"user"`
}

type TokenPairResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
