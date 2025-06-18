package neo4jOperations

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go-api/internal/database"
	"go-api/internal/models"
	"strings"
)

func CreateWebResource(wr models.WebResource) (*string, error) {

	params := map[string]interface{}{
		"type":     strings.ToLower(wr.Type),
		"value":    strings.ToLower(wr.Value),
		"labels":   wr.Labels,
		"disabled": wr.Disabled,
	}

	keys := make([]string, 0, len(params))
	for key := range params {
		keys = append(keys, "wr."+key+"=$"+key)
	}
	strParams := "SET " + strings.Join(keys, ", ")

	query := fmt.Sprintf("CREATE (wr:WebResource) %s RETURN id(wr) AS id", strParams)

	driver, err := database.NewNeo4jDriver()
	if err != nil {
		return nil, err
	}
	defer database.CloseNeo4jDriver(driver)

	session := driver.NewSession(context.Background(), neo4j.SessionConfig{})
	defer database.CloseNeo4jSession(session)

	result, err := session.Run(context.Background(), query, params)
	if err != nil {
		return nil, fmt.Errorf("query execution failed: %w", err)
	}
	if result.Err() != nil {
		return nil, fmt.Errorf("error iterating results: %w", result.Err())
	}

	if result.Next(context.Background()) {
		record := result.Record()
		id, _ := record.Get("id")
		idStr := fmt.Sprintf("%v", id)
		return &idStr, nil
	} else {
		return nil, fmt.Errorf("query execution failed: %w", err)
	}
}
