package models

type LoginResponse struct {
	Token                 string    `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	TokenExpiresAt        int64     `json:"token_expires_at"`
	RefreshToken          string    `json:"refresh_token"`
	RefreshTokenExpiresAt int64     `json:"refresh_token_expires_at"`
	Role                  UserRoles `json:"role"`
}
