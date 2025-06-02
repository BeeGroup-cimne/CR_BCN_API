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

func UpdateUser(user models.UpdateUserRequest, value interface{}) error {
	var params = map[string]interface{}{}
	var id int
	var email string

	switch value.(type) {
	case int:
		id = value.(int)
		break
	case string:
		email = strings.ToLower(value.(string))
		break
	default:
		return fmt.Errorf("value must be either an email(string) or id(int)")
	}

	if user.Email != nil {
		params["email"] = strings.ToLower(*user.Email)
	}
	if user.Password != nil {
		hash, err := password.HashPassword(*user.Password)
		if err != nil {
			return err
		}
		params["password"] = hash
		params["password_reset_token_hash"] = nil
		params["password_reset_token_expires_at"] = nil
	}
	if user.Role != nil {
		params["role"] = *user.Role
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
	if user.Entity != nil {
		params["entity"] = *user.Entity
	}
	if user.Department != nil {
		params["department"] = *user.Department
	}
	if user.Name != nil {
		params["name"] = *user.Name
	}
	if user.PasswordResetTokenHash != nil {
		params["password_reset_token_hash"] = *user.PasswordResetTokenHash
	}
	if user.PasswordResetTokenExpiresAt != nil {
		expiresAtFormatted := user.PasswordResetTokenExpiresAt.UTC().Format(time.RFC3339)
		params["password_reset_token_expires_at"] = expiresAtFormatted
	}

	keys := make([]string, 0, len(params))
	for key := range params {
		keys = append(keys, "u."+key+"=$"+key)
	}
	strParams := "SET " + strings.Join(keys, ", ")
	params["user_id"] = id
	params["user_email"] = email

	query := fmt.Sprintf("MATCH (u:APIUser) WHERE id(u)=$user_id OR u.email=$user_email %s", strParams)

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
