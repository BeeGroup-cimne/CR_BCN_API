package neo4jOperations

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go-api/internal/database"
	_ "strconv"
)

func GetFilteredClimateShelterDetailsById(id string) (map[string]interface{}, error) {

	query := "MATCH (n:ns1__ClimateSheltersBuildingSpace) where n.ns1__climateShelterId = $id " +
		"WITH replace(replace(n.ns1__schedule, '&apos;', '\\''), '\\'', '\\\"') AS cleanedJson, n " +
		"WITH apoc.convert.fromJsonList(cleanedJson) AS scheduleJson, n " +
		"UNWIND scheduleJson AS outerList " +
		"UNWIND outerList AS schedule " +
		"WITH schedule, n "

	query += "RETURN n { .*, schedule: schedule} AS result"

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
		record := result.Record()

		value, ok := record.Get("result")
		if ok {
			strValue, ok := value.(map[string]interface{})
			if ok {
				finalResult = strValue
			}
		}

	}

	if err := result.Err(); err != nil {
		return nil, err
	}

	return finalResult, nil
}
