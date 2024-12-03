package types

type UserResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type TokenPairResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

type ValidateResponse struct {
	Valid bool         `json:"valid"`
	User  UserResponse `json:"user"`
}

type ValidateTokenResponse struct {
	Subject string `json:"subject"`
}
