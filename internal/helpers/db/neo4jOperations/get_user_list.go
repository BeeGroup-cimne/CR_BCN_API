package neo4jOperations

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go-api/internal/database"
	"go-api/internal/helpers/transform"
	"go-api/internal/models"
)

func GetUserList() ([]*models.User, error) {
	query := "MATCH (u:APIUser) RETURN u {.*, id: id(u)} AS user"

	driver, err := database.NewNeo4jDriver()
	if err != nil {
		return nil, err
	}
	defer database.CloseNeo4jDriver(driver)

	session := driver.NewSession(context.Background(), neo4j.SessionConfig{})
	defer database.CloseNeo4jSession(session)

	result, err := session.Run(context.Background(), query, nil)
	if err != nil {
		return nil, err
	}

	var usersList []*models.User

	for result.Next(context.Background()) {
		record := result.Record()
		userMap, ok := record.Get("user")
		if !ok {
			continue
		}

		data, ok := userMap.(map[string]interface{})
		if !ok {
			continue
		}

		user := &models.User{
			Id: int(data["id"].(int64)),
			//Password:                     data["password"].(string),
			Email:                        data["email"].(string),
			Department:                   transform.ConvertToNullableString(data["department"]),
			Name:                         transform.ConvertToNullableString(data["name"]),
			Entity:                       transform.ConvertToNullableString(data["entity"]),
			CreationDate:                 data["creation_date"].(string),
			Role:                         models.UserRoles(data["role"].(string)),
			LastConnectionDate:           transform.ConvertToNullableString(data["last_connection_date"]),
			AccumulatedHourlyConnections: transform.ConvertToNullableInt(data["accumulated_hourly_connections"]),
			MaxHourlyConnections:         transform.ConvertToNullableInt(data["max_hourly_connections"]),
			AccumulatedDailyConnections:  transform.ConvertToNullableInt(data["accumulated_daily_connections"]),
			MaxDailyConnections:          transform.ConvertToNullableInt(data["max_daily_connections"]),
			AccumulatedTotalConnections:  transform.ConvertToNullableInt(data["accumulated_total_connections"]),
			Blocked:                      transform.ConvertToNullableBool(data["blocked"]),
			Token:                        transform.ConvertToNullableString(data["token"]),
			RefreshToken:                 transform.ConvertToNullableString(data["refresh_token"]),
			ExpirationDate:               transform.ConvertToNullableString(data["expiration_date"]),
			StartDate:                    transform.ConvertToNullableString(data["start_date"]),
		}

		usersList = append(usersList, user)
	}

	if err = result.Err(); err != nil {
		return nil, err
	}

	return usersList, nil
}
