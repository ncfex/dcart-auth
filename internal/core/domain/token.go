package domain

type Token string

type TokenPair struct {
	AccessToken  Token `json:"access_token"`
	RefreshToken Token `json:"refresh_token"`
}
