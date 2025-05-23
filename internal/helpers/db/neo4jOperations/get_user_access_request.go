package neo4jOperations

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go-api/internal/database"
	"go-api/internal/models"
)

func GetUserAccessRequest(value interface{}) (*models.UserAccessRequest, error) {

	query := "MATCH (u:UserAccessRequest) WHERE id(u)=$id RETURN u {.*, id: id(u)} AS userAccessRequest"
	var id int64

	switch value.(type) {
	case int64:
		id = value.(int64)
		break
	default:
		return nil, fmt.Errorf("value must be either an email(string)")
	}

	params := map[string]interface{}{
		"id": id,
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
			Id:         data["id"].(int64),
			Password:   data["password"].(string),
			Email:      data["email"].(string),
			Entity:     data["entity"].(string),
			Department: data["department"].(string),
			Name:       data["name"].(string),
			Status:     models.UserAccessRequestStatus(data["status"].(string)),
		}, nil
	}

	if err = result.Err(); err != nil {
		return nil, err
	}

	return nil, nil
}
