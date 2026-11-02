package neo4jOperations

import (
	"context"
	"log"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go-api/internal/database"
)

func GetFilteredClimateShelter(date string, characteristics map[string]bool, buildingType map[string]bool) ([]string, int, int, error) {

	selectedCharacteristics := selectedKeys(characteristics)
	selectedBuildingType := selectedKeys(buildingType)

	query := "MATCH (n:ns1__ClimateSheltersBuildingSpace) " +
		climateShelterFilterConditions(selectedCharacteristics, selectedBuildingType) +
		"RETURN n.ns1__climateShelterId AS id, n.ns1__schedule AS schedule"

	var t time.Time
	if date != "" {
		var err error
		t, err = time.Parse(time.RFC3339Nano, date)
		if err != nil {
			return nil, 0, 0, err
		}
	}

	driver, err := database.NewNeo4jDriver()
	if err != nil {
		return nil, 0, 0, err
	}
	defer database.CloseNeo4jDriver(driver)

	session := driver.NewSession(context.Background(), neo4j.SessionConfig{})
	defer database.CloseNeo4jSession(session)

	params := map[string]any{
		"selectedCharacteristics": selectedCharacteristics,
		"selectedBuildingType":    selectedBuildingType,
	}

	result, err := session.Run(context.Background(), query, params)
	if err != nil {
		return nil, 0, 0, err
	}

	var shelters []string
	for result.Next(context.Background()) {
		record := result.Record()

		val, _ := record.Get("id")
		id, ok := val.(string)
		if !ok {
			continue
		}

		if date != "" {
			raw, _ := record.Get("schedule")
			rawStr, ok := raw.(string)
			if !ok {
				continue
			}
			open, err := shelterIsOpenAt(rawStr, t)
			if err != nil {
				log.Printf("refugio %s: horario no válido: %v", id, err)
				continue
			}
			if !open {
				continue
			}
		}
		shelters = append(shelters, id)
	}
	if err := result.Err(); err != nil {
		return nil, 0, 0, err
	}

	totalResult, err := session.Run(context.Background(),
		"MATCH (s:ns1__ClimateSheltersBuildingSpace) RETURN COUNT(s) AS total_n_shelters", nil)
	if err != nil {
		return nil, 0, 0, err
	}
	totalNShelters := 0
	if totalResult.Next(context.Background()) {
		if val, ok := totalResult.Record().Get("total_n_shelters"); ok {
			if count, ok := val.(int64); ok {
				totalNShelters = int(count)
			}
		}
	}
	if err := totalResult.Err(); err != nil {
		return nil, 0, 0, err
	}

	return shelters, len(shelters), totalNShelters, nil
}
