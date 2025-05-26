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

func CreateUser(user models.CreateUserRequest) (*string, error) {
	hash, err := password.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}
	params := map[string]interface{}{
		"email":                          strings.ToLower(user.Email),
		"password":                       hash,
		"creation_date":                  time.Now().UTC().Format(time.RFC3339),
		"role":                           user.Role,
		"accumulated_hourly_connections": 0,
		"accumulated_daily_connections":  0,
		"accumulated_total_connections":  0,
	}
	if user.Entity != nil {
		params["entity"] = *user.Entity
	}
	if user.Department != nil {
		params["department"] = *user.Department
	}
	if user.Name != nil {
		params["name"] = *user.Name
	}
	if user.MaxHourlyConnections != nil {
		params["max_hourly_connections"] = *user.MaxHourlyConnections
	}
	if user.MaxDailyConnections != nil {
		params["max_daily_connections"] = *user.MaxDailyConnections
	}
	if user.Blocked != nil {
		params["blocked"] = *user.Blocked
	}
	if user.ExpirationDate != nil {
		params["expiration_date"] = *user.ExpirationDate
	}
	if user.StartDate != nil {
		params["start_date"] = *user.StartDate
	}

	keys := make([]string, 0, len(params))
	for key := range params {
		keys = append(keys, "u."+key+"=$"+key)
	}
	strParams := "SET " + strings.Join(keys, ", ")

	query := fmt.Sprintf("CREATE (u:APIUser) %s RETURN id(u) AS id", strParams)

	driver, err := database.NewNeo4jDriver()
	if err != nil {
		return nil, err
	}
	defer database.CloseNeo4jDriver(driver)

	session := driver.NewSession(context.Background(), neo4j.SessionConfig{})
	defer database.CloseNeo4jSession(session)

	result, err := session.Run(context.Background(), query, params)
	if err != nil {
		return nil, fmt.Errorf("query execution failed: %w", err)
	}
	if result.Err() != nil {
		return nil, fmt.Errorf("error iterating results: %w", result.Err())
	}

	if result.Next(context.Background()) {
		record := result.Record()
		id, _ := record.Get("id")
		idStr := fmt.Sprintf("%v", id)
		return &idStr, nil
	} else {
		return nil, fmt.Errorf("query execution failed: %w", err)
	}
}
