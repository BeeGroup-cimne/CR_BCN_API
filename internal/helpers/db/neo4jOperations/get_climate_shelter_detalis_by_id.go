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

	query += "RETURN { " +
		"climateShelterId: n.ns1__climateShelterId, " +
		"hasAccessibleWc: n.ns1__has_accessible_wc, " +
		"hasDrinkingFountain: n.ns1__has_drinking_fountain, " +
		"hasDrinkingWater: n.ns1__has_drinking_water, " +
		"hasKidsArea: n.ns1__has_kids_area, " +
		"hasMedicalService: n.ns1__has_medical_service, " +
		"hasNursingService: n.ns1__has_nursing_service, " +
		"hasPetsAllowed: n.ns1__has_pets_allowed, " +
		"hasRestroom: n.ns1__has_restroom, " +
		"hasSummerClimateShelter: n.ns1__has_summer_climate_shelter, " +
		"hasWheelchairSpace: n.ns1__has_wheelchair_space, " +
		"hasWifi: n.ns1__has_wifi, " +
		"hasWinterClimateShelter: n.ns1__has_winter_climate_shelter, " +
		"isAccessible: n.ns1__is_accessible, " +
		"isAccessiblePhysical: n.ns1__is_accessible_physical, " +
		"name: n.ns1__name, " +
		"streetName: n.ns1__street_name, " +
		"streetNumber: n.ns1__street_number, " +
		"type: n.ns1__type, " +
		"schedule: schedule " +
		"} AS result"

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
