package dbOperations

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go-api/internal/database"
	"go-api/internal/helpers/transform"
	"go-api/internal/models"
	"strings"
)

type Neo4jResponse struct {
	User models.User `json:"user"`
}

func FetchNeo4jUser(email string) (*models.User, error) {

	query := "MATCH (u:APIUser) WHERE u.email=$email RETURN u {.*, id: id(u)} AS user LIMIT 1"
	params := map[string]interface{}{
		"email": strings.ToLower(email),
	}

	driver, err := database.NewNeo4jDriver()
	if err != nil {
		return nil, err
	}
	defer database.CloseNeo4jDriver(driver)

	session := driver.NewSession(context.Background(), neo4j.SessionConfig{})
	defer database.CloseNeo4jSession(session)

	result, err := session.Run(context.Background(), query, params)
	if err != nil {
		return nil, err
	}

	// Extract result and map to struct
	if result.Next(context.Background()) {
		record := result.Record()
		userMap, _ := record.Get("user")
		data := userMap.(map[string]interface{})

		return &models.User{
			Id:                           int(data["id"].(int64)),
			Password:                     data["password"].(string),
			Email:                        data["email"].(string),
			CreationDate:                 data["creation_date"].(string),
			Role:                         models.UserRoles(data["role"].(string)),
			LastConnectionDate:           transform.ConvertToNullableString(data["last_connection_date"]),
			AccumulatedHourlyConnections: transform.ConvertToNullableInt(data["accumulated_hourly_connections"]),
			MaxHourlyConnections:         transform.ConvertToNullableInt(data["max_hourly_connections"]),
			AccumulatedDailyConnections:  transform.ConvertToNullableInt(data["accumulated_daily_connections"]),
			MaxDailyConnections:          transform.ConvertToNullableInt(data["max_daily_connections"]),
			Blocked:                      transform.ConvertToNullableBool(data["blocked"]),
			Token:                        transform.ConvertToNullableString(data["token"]),
			RefreshToken:                 transform.ConvertToNullableString(data["refresh_token"]),
			ExpirationDate:               transform.ConvertToNullableString(data["expiration_date"]),
			StartDate:                    transform.ConvertToNullableString(data["start_date"]),
		}, nil
	}

	if err = result.Err(); err != nil {
		return nil, err
	}

	return nil, fmt.Errorf("dbOperations not found")
}
