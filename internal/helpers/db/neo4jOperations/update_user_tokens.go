package neo4jOperations

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go-api/internal/database"
)

func UpdateUserTokens(id int, token string, refreshToken string) error {
	query := "MATCH (u) WHERE id(u)=$id SET u.token=$token, u.refresh_token=$refresh_token;"
	params := map[string]interface{}{
		"id":            id,
		"token":         token,
		"refresh_token": refreshToken,
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
