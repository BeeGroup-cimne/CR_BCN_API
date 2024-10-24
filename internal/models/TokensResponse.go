package models

type TokensResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}
