package neo4jOperations

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go-api/internal/database"
	"go-api/internal/helpers/transform"
	"go-api/internal/helpers/utils"
	"go-api/internal/models"
	"strings"
)

func GetWebResourceListByType(value string) ([]*models.WebResource, error) {
	query := "MATCH (wr:WebResource) WHERE wr.type=$value RETURN wr {.*, id: id(wr)} AS wr"
	params := map[string]interface{}{
		"value": strings.ToLower(value),
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

	var wrList []*models.WebResource

	for result.Next(context.Background()) {
		record := result.Record()
		wrMap, ok := record.Get("wr")
		if !ok {
			continue
		}

		data, ok := wrMap.(map[string]interface{})
		if !ok {
			continue
		}

		wr := &models.WebResource{
			Id:       transform.ConvertToNullableInt(data["id"]),
			Type:     data["type"].(string),
			Value:    data["value"].(string),
			Labels:   utils.ToStringArray(data["labels"]),
			Disabled: data["disabled"].(bool),
		}

		wrList = append(wrList, wr)
	}

	if err = result.Err(); err != nil {
		return nil, err
	}

	return wrList, nil
}
