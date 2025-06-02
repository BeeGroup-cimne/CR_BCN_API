package models

import "time"

type UpdateUserRequest struct {
	Email                       *string    `json:"email"`
	Password                    *string    `json:"password"`
	Role                        *UserRoles `json:"role"`
	MaxHourlyConnections        *int       `json:"max_hourly_connections"`
	MaxDailyConnections         *int       `json:"max_daily_connections"`
	Blocked                     *bool      `json:"blocked"`
	ExpirationDate              *string    `json:"expiration_date"`
	StartDate                   *string    `json:"start_date"`
	Entity                      *string    `json:"entity"`
	Department                  *string    `json:"department"`
	Name                        *string    `json:"name"`
	PasswordResetTokenHash      *string    `json:"password_reset_token_hash"`
	PasswordResetTokenExpiresAt *time.Time `json:"password_reset_token_expires_at"`
}
