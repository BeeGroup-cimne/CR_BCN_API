package neo4jOperations

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go-api/internal/database"
	"go-api/internal/helpers/security/password"
	"go-api/internal/models"
	"strings"
	"time"
)

func CreateUser(user models.CreateUserRequest) error {
	hash, err := password.HashPassword(user.Password)
	if err != nil {
		return err
	}
	params := map[string]interface{}{
		"email":         strings.ToLower(user.Email),
		"password":      hash,
		"creation_date": time.Now().UTC().Format(time.RFC3339),
		"role":          user.Role,
	}
	var strParams string
	if user.MaxHourlyConnections != nil {
		strParams += ", max_hourly_connections:$max_hourly_connections"
		params["max_hourly_connections"] = *user.MaxHourlyConnections
	}
	if user.MaxDailyConnections != nil {
		strParams += ", max_daily_connections:$max_daily_connections"
		params["max_daily_connections"] = *user.MaxDailyConnections
	}
	if user.Blocked != nil {
		strParams += ", blocked:$blocked"
		params["blocked"] = *user.Blocked
	}
	if user.ExpirationDate != nil {
		strParams += ", expiration_date:$expiration_date"
		params["expiration_date"] = *user.ExpirationDate
	}
	if user.StartDate != nil {
		strParams += ", start_date:$start_date"
		params["start_date"] = *user.StartDate
	}
	query := fmt.Sprintf("CREATE (u:APIUser{email:$email, password:$password, creation_date:$creation_date, role:$role, accumulated_hourly_connections: 0, accumulated_daily_connections: 0, accumulated_total_connections: 0%s})", strParams)

	driver, err := database.NewNeo4jDriver()
	if err != nil {
		return err
	}
	defer database.CloseNeo4jDriver(driver)

	session := driver.NewSession(context.Background(), neo4j.SessionConfig{})
	defer database.CloseNeo4jSession(session)

	result, err := session.Run(context.Background(), query, params)
	if err != nil {
		return fmt.Errorf("query execution failed: %w", err)
	}
	if result.Err() != nil {
		return fmt.Errorf("error iterating results: %w", result.Err())
	}
	return nil
}
