package neo4jOperations

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go-api/internal/database"
	"go-api/internal/models"
	userAccessRequestStatus "go-api/internal/models/enums/UserAccesRequestStatus"
)

func GetUserAccessRequestList() ([]*models.UserAccessRequest, error) {
	query := "MATCH (u:UserAccessRequest) WHERE u.status=$status RETURN u {.*, id: id(u)} AS userAccessRequest"

	params := map[string]interface{}{
		"status": userAccessRequestStatus.Unprocessed,
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

	var userRequests []*models.UserAccessRequest

	for result.Next(context.Background()) {
		record := result.Record()
		userMap, ok := record.Get("userAccessRequest")
		if !ok {
			continue
		}

		data, ok := userMap.(map[string]interface{})
		if !ok {
			continue
		}

		user := &models.UserAccessRequest{
			Id:         data["id"].(int64),
			Password:   data["password"].(string),
			Email:      data["email"].(string),
			Entity:     data["entity"].(string),
			Department: data["department"].(string),
			Name:       data["name"].(string),
			Status:     models.UserAccessRequestStatus(data["status"].(string)),
		}

		userRequests = append(userRequests, user)
	}

	if err = result.Err(); err != nil {
		return nil, err
	}

	return userRequests, nil
}
