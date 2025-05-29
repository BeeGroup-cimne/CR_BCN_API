package models

type User struct {
	Id                           int       `json:"id"`
	Email                        string    `json:"email"`
	CreationDate                 string    `json:"creation_date"`
	Role                         UserRoles `json:"role"`
	Name                         *string   `json:"name"`
	Department                   *string   `json:"department"`
	Entity                       *string   `json:"entity"`
	LastConnectionDate           *string   `json:"last_connection_date"`
	AccumulatedHourlyConnections *int      `json:"accumulated_hourly_connections"`
	MaxHourlyConnections         *int      `json:"max_hourly_connections"`
	AccumulatedDailyConnections  *int      `json:"accumulated_daily_connections"`
	MaxDailyConnections          *int      `json:"max_daily_connections"`
	AccumulatedTotalConnections  *int      `json:"accumulated_total_connections"`
	Blocked                      *bool     `json:"blocked"`
	Token                        *string   `json:"token"`
	RefreshToken                 *string   `json:"refresh_token"`
	ExpirationDate               *string   `json:"expiration_date"`
	StartDate                    *string   `json:"start_date"`
	Password                     string    `json:"password"`
}
