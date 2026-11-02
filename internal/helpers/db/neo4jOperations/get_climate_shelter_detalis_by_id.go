package neo4jOperations

import (
	"context"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go-api/internal/database"
)

func GetFilteredClimateShelterDetailsById(id string) (map[string]interface{}, error) {

	query := "MATCH (n:ns1__ClimateSheltersBuildingSpace) " +
		"WHERE n.ns1__climateShelterId = $id " +
		"RETURN properties(n) AS props"

	driver, err := database.NewNeo4jDriver()
	if err != nil {
		return nil, err
	}
	defer database.CloseNeo4jDriver(driver)

	session := driver.NewSession(context.Background(), neo4j.SessionConfig{})
	defer database.CloseNeo4jSession(session)

	params := map[string]any{
		"id": id,
	}

	result, err := session.Run(context.Background(), query, params)
	if err != nil {
		return nil, err
	}

	finalResult := map[string]interface{}{}
	for result.Next(context.Background()) {
		value, ok := result.Record().Get("props")
		if !ok {
			continue
		}
		props, ok := value.(map[string]interface{})
		if !ok {
			continue
		}

		raw, _ := props["ns1__schedule"].(string)
		schedules, err := parseScheduleList(raw)
		if err != nil {
			return nil, fmt.Errorf("horario no válido en el refugio %s: %w", id, err)
		}

		for _, s := range schedules {
			row := make(map[string]interface{}, len(props)+1)
			for k, v := range props {
				row[k] = v
			}
			row["schedule"] = s
			finalResult = row
		}
	}

	if err := result.Err(); err != nil {
		return nil, err
	}

	return finalResult, nil
}
