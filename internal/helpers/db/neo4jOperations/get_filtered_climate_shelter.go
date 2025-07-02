package neo4jOperations

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go-api/internal/database"
	_ "strconv"
	"time"
)

func GetFilteredClimateShelter(date string, characteristics map[string]bool, buildingType map[string]bool) ([]string, int, int, error) {

	//query := "MATCH (n:ns1__ClimateSheltersBuildingSpace) "
	query := "CALL{MATCH (n:ns1__ClimateSheltersBuildingSpace) "

	selectedCharacteristics := []string{}
	for key, value := range characteristics {
		if value {
			selectedCharacteristics = append(selectedCharacteristics, key)
		}
	}
	selectedBuildingType := []string{}
	for key, value := range buildingType {
		if value {
			selectedBuildingType = append(selectedBuildingType, key)
		}
	}

	if len(selectedCharacteristics) != 0 {
		query += "WHERE ALL(attr IN $selectedCharacteristics WHERE n[attr] = \"True\") "
	}

	if len(selectedBuildingType) != 0 {
		query += "and  n.ns1__type IN $selectedBuildingType "
	}

	if date != "" {
		query +=
			"WITH replace(replace(n.ns1__schedule, '&apos;', '\\''), '\\'', '\\\"') AS cleanedJson, " +
				"date($currentDate) AS currentDate, $weekday AS weekday, $hour AS hour, n " +
				"WITH apoc.convert.fromJsonList(cleanedJson) AS scheduleJson, n, currentDate, weekday, hour " +
				"UNWIND scheduleJson AS outerList " +
				"UNWIND outerList AS schedule " +
				"WITH schedule, n, currentDate, weekday, hour " +
				"UNWIND (CASE WHEN size(schedule.periods) = 0 THEN [NULL] ELSE schedule.periods END) AS period " +
				"WITH schedule, n, weekday, hour, period, currentDate " +
				"WHERE period IS NULL " +
				"OR (( period.month_start <= period.month_end AND currentDate.month >= period.month_start AND currentDate.month <= period.month_end AND currentDate.day >= period.day_start AND currentDate.day <= period.day_end) " +
				"OR ( period.month_start > period.month_end AND ( currentDate.month >= period.month_start OR currentDate.month <= period.month_end) AND currentDate.day >= period.day_start AND currentDate.day <= period.day_end)) " +
				"UNWIND schedule.weekly_profile AS profile WITH profile, n, weekday, hour WHERE weekday IN profile.days " +
				"UNWIND profile.opening_hours AS opening WITH opening, n, hour " +
				"WHERE ((opening.hour_start <= opening.hour_end AND opening.hour_start <= hour AND opening.hour_end >= hour) " +
				"OR(opening.hour_start > opening.hour_end AND (hour >= opening.hour_start OR hour <= opening.hour_end))) "
	}

	//query += "RETURN n.ns1__climateShelterId AS result"
	query += "RETURN COLLECT(n.ns1__climateShelterId) AS shelters, COUNT(n) AS n_shelters}"
	query += "CALL { MATCH (s:ns1__ClimateSheltersBuildingSpace) RETURN COUNT(s) AS total_n_shelters }"
	query += "RETURN shelters, n_shelters, total_n_shelters"

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

	if date != "" {
		layout := time.RFC3339Nano
		t, err := time.Parse(layout, date)
		if err != nil {
			return nil, 0, 0, err
		}
		currentDate := t.Format("2006-01-02")
		weekday := t.Weekday().String()
		hour := t.Format("15:04")

		params["currentDate"] = currentDate
		params["weekday"] = weekday
		params["hour"] = hour
	}

	result, err := session.Run(context.Background(), query, params)
	if err != nil {
		return nil, 0, 0, err
	}
	var shelters []string
	var nShelters, totalNShelters int

	for result.Next(context.Background()) {
		record := result.Record()

		if val, ok := record.Get("shelters"); ok {
			if list, ok := val.([]interface{}); ok {
				for _, id := range list {
					if strId, ok := id.(string); ok {
						shelters = append(shelters, strId)
					}
				}
			}
		}

		if val, ok := record.Get("n_shelters"); ok {
			if count, ok := val.(int64); ok {
				nShelters = int(count)
			}
		}

		if val, ok := record.Get("total_n_shelters"); ok {
			if count, ok := val.(int64); ok {
				totalNShelters = int(count)
			}
		}
	}

	if err := result.Err(); err != nil {
		return nil, 0, 0, err
	}

	return shelters, nShelters, totalNShelters, nil
}
