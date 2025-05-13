package neo4jOperations

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go-api/internal/database"
	"go-api/internal/models"
	"strings"
)

func GetUserAccessRequest(value interface{}) (*models.UserAccessRequest, error) {

	query := "MATCH (u:UserAccessRequest) WHERE u.email=$email RETURN u {.*, id: id(u)} AS userAccessRequest LIMIT 1"
	var email string

	switch value.(type) {
	case string:
		email = strings.ToLower(value.(string))
		break
	default:
		return nil, fmt.Errorf("value must be either an email(string)")
	}

	params := map[string]interface{}{
		"email": email,
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

	if result.Next(context.Background()) {
		record := result.Record()
		userMap, _ := record.Get("userAccessRequest")
		data := userMap.(map[string]interface{})

		return &models.UserAccessRequest{
			Password:   data["password"].(string),
			Email:      data["email"].(string),
			Entity:     data["entity"].(string),
			Department: data["department"].(string),
			Name:       data["name"].(string),
		}, nil
	}

	if err = result.Err(); err != nil {
		return nil, err
	}

	return nil, nil
}
