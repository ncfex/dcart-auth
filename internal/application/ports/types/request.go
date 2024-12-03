package types

type RegisterRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type TokenRequest struct {
	Token string `json:"token" validate:"required"`
}

type CreateTokenParams struct {
	UserID string `json:"user_id" validate:"required"`
}
