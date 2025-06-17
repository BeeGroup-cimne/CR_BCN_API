package neo4jOperations

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go-api/internal/database"
	"go-api/internal/models"
	"strings"
)

func UpdateWebResource(id int, wr models.WebResourceUpdate) error {
	params := map[string]interface{}{
		"id": id,
	}

	if wr.Type != nil {
		params["type"] = strings.ToLower(*wr.Type)
	}
	if wr.Value != nil {
		params["value"] = strings.ToLower(*wr.Value)
	}
	if wr.Labels != nil {
		params["labels"] = *wr.Labels
	}
	if wr.Disabled != nil {
		params["disabled"] = *wr.Disabled
	}

	keys := make([]string, 0, len(params))
	for key := range params {
		keys = append(keys, "wr."+key+"=$"+key)
	}
	strParams := "SET " + strings.Join(keys, ", ")
	query := fmt.Sprintf("MATCH (wr:WebResource) WHERE id(wr)=$id %s", strParams)

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
