package neo4jOperations

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go-api/internal/database"
)

func DeleteUser(id int) error {
	query := "MATCH (u:APIUser) WHERE id(u)=$id DETACH DELETE u"
	params := map[string]interface{}{
		"id": id,
	}

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
