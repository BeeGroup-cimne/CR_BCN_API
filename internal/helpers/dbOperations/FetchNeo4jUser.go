package dbOperations

import (
	"context"
	"fmt"
	"go-api/internal/database"
	"go-api/internal/helpers/transform"
	"go-api/internal/models"
	"strings"
)

func FetchNeo4jUser(email string, password string) (*models.User, error) {
	query := "MATCH (n:APIUser) WHERE u.email=$email AND u.password=$password RETURN n {.*, id: id(n)} LIMIT 1"
	params := map[string]interface{}{
		"email":    strings.ToLower(email),
		"password": password,
	}

	result, err := database.RunNeo4jQuery(query, params)

	if err != nil {
		return nil, err
	}

	// Extract result and map to struct
	if result.Next(context.Background()) {
		record := result.Record()

		// Extract fields from the record
		id, _ := record.Get("id")
		password, _ := record.Get("password")
		email, _ := record.Get("email")
		creationDate, _ := record.Get("creation_date")
		lastConnectionDate, _ := record.Get("last_connection_date")
		accumulatedHourlyConnections, _ := record.Get("accumulated_hourly_connections")
		maxHourlyConnections, _ := record.Get("max_hourly_connections")
		accumulatedDailyConnections, _ := record.Get("accumulated_daily_connections")
		maxDailyConnections, _ := record.Get("max_daily_connections")
		blocked, _ := record.Get("blocked")
		token, _ := record.Get("token")
		refreshToken, _ := record.Get("refresh_token")
		expirationDate, _ := record.Get("expiration_date")
		startDate, _ := record.Get("start_date")
		role, _ := record.Get("role")

		// Create a User struct and populate it
		user := &models.User{
			Id:                           id.(int),
			Password:                     password.(string),
			Email:                        email.(string),
			CreationDate:                 creationDate.(string),
			LastConnectionDate:           transform.ConvertToNullableString(lastConnectionDate.(string)),
			AccumulatedHourlyConnections: transform.ConvertToNullableInt(accumulatedHourlyConnections.(int)),
			MaxHourlyConnections:         transform.ConvertToNullableInt(maxHourlyConnections.(int)),
			AccumulatedDailyConnections:  transform.ConvertToNullableInt(accumulatedDailyConnections.(int)),
			MaxDailyConnections:          transform.ConvertToNullableInt(maxDailyConnections.(int)),
			Blocked:                      transform.ConvertToNullableBool(blocked.(bool)),
			Token:                        transform.ConvertToNullableString(token.(string)),
			RefreshToken:                 transform.ConvertToNullableString(refreshToken.(string)),
			ExpirationDate:               transform.ConvertToNullableString(expirationDate.(string)),
			StartDate:                    transform.ConvertToNullableString(startDate.(string)),
			Role:                         models.UserRoles(role.(string)),
		}

		return user, nil
	}

	if err = result.Err(); err != nil {
		return nil, err
	}

	return nil, fmt.Errorf("dbOperations not found")
}
