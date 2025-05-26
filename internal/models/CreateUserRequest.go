package models

type CreateUserRequest struct {
	Email                string    `json:"email"`
	Password             string    `json:"password"`
	Role                 UserRoles `json:"role"`
	Name                 *string   `json:"name"`
	Department           *string   `json:"department"`
	Entity               *string   `json:"entity"`
	MaxHourlyConnections *int      `json:"max_hourly_connections"`
	MaxDailyConnections  *int      `json:"max_daily_connections"`
	Blocked              *bool     `json:"blocked"`
	ExpirationDate       *string   `json:"expiration_date"`
	StartDate            *string   `json:"start_date"`
}
