package neo4jOperations

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go-api/internal/database"
	"go-api/internal/models"
)

func UpdateUserAccessRequestStatus(user models.UpdateStatusAccessRequest) error {
	query := "MATCH (u) WHERE u.email=$email SET u.status=$status;"
	params := map[string]interface{}{
		"email":  user.Email,
		"status": user.Status,
	}
	driver, err := database.NewNeo4jDriver()
	if err != nil {
		return err
	}
	defer database.CloseNeo4jDriver(driver)

	session := driver.NewSession(context.Background(), neo4j.SessionConfig{})
	defer database.CloseNeo4jSession(session)

	_, err = session.Run(context.Background(), query, params)
	if err != nil {
		return fmt.Errorf("query execution failed: %w", err)
	}
	return nil
}
