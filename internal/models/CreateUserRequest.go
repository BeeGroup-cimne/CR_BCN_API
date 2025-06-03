package models

import "time"

type CreateUserRequest struct {
	Email                       string     `json:"email"`
	Role                        UserRoles  `json:"role"`
	Password                    *string    `json:"password"`
	Name                        *string    `json:"name"`
	Department                  *string    `json:"department"`
	Entity                      *string    `json:"entity"`
	MaxHourlyConnections        *int       `json:"max_hourly_connections"`
	MaxDailyConnections         *int       `json:"max_daily_connections"`
	Blocked                     *bool      `json:"blocked"`
	ExpirationDate              *string    `json:"expiration_date"`
	StartDate                   *string    `json:"start_date"`
	PasswordResetTokenHash      *string    `json:"password_reset_token_hash"`
	PasswordResetTokenExpiresAt *time.Time `json:"password_reset_token_expires_at"`
}
