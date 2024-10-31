package models

type CreateUserRequest struct {
	Email                string    `json:"email"`
	Password             string    `json:"password"`
	Role                 UserRoles `json:"role"`
	MaxHourlyConnections *int      `json:"max_hourly_connections"`
	MaxDailyConnections  *int      `json:"max_daily_connections"`
	Blocked              *bool     `json:"blocked"`
	ExpirationDate       *string   `json:"expiration_date"`
	StartDate            *string   `json:"start_date"`
}
