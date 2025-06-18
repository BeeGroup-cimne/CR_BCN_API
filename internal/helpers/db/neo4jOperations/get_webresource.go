package neo4jOperations

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go-api/internal/database"
	"go-api/internal/helpers/transform"
	"go-api/internal/helpers/utils"
	"go-api/internal/models"
)

func GetWebResource(id int) (*models.WebResource, error) {

	query := "MATCH (wr:WebResource) WHERE id(wr)=$id RETURN wr {.*, id: id(wr)} AS wr LIMIT 1"

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
		userMap, _ := record.Get("wr")
		data := userMap.(map[string]interface{})

		return &models.WebResource{
			Id:       transform.ConvertToNullableInt(data["id"]),
			Type:     data["type"].(string),
			Value:    data["value"].(string),
			Labels:   utils.ToStringArray(data["labels"]),
			Disabled: data["disabled"].(bool),
		}, nil
	}

	if err = result.Err(); err != nil {
		return nil, err
	}

	return nil, fmt.Errorf("WebResource does not exist")
}
