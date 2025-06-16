package neo4jOperations

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go-api/internal/database"
	_ "strconv"
	"time"
)

func GetFilteredClimateShelter(date string, characteristics map[string]bool, buildingType map[string]bool) ([]string, error) {

	query := "MATCH (n:ns1__ClimateSheltersBuildingSpace) "

	selectedCharacteristics := []string{}
	for key, value := range characteristics {
		if value {
			selectedCharacteristics = append(selectedCharacteristics, "ns1__has_"+key)
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
		query += "WITH replace(replace(n.ns1__schedule, '&apos;', '\\''), '\\'', '\\\"') AS cleanedJson, " +
			"date($currentDate) AS currentDate, $weekday AS weekday, $hour AS hour, n " +
			"WITH apoc.convert.fromJsonList(cleanedJson) AS scheduleJson, n, currentDate, weekday, hour " +
			"UNWIND scheduleJson AS outerList " +
			"UNWIND outerList AS schedule " +
			"WITH schedule, n, currentDate, weekday, hour " +
			"UNWIND (CASE WHEN size(schedule.periods) = 0 THEN [NULL] ELSE schedule.periods END) AS period " +
			"WITH schedule, n, weekday, hour, period, currentDate " +
			"WHERE period IS NULL OR ( currentDate.month >= period.month_start AND currentDate.month <= period.month_end AND currentDate.day >= period.day_start AND currentDate.day <= period.day_end) " +
			"UNWIND schedule.weekly_profile AS profile WITH profile, n, weekday, hour WHERE weekday IN profile.days " +
			"UNWIND profile.opening_hours AS opening WITH opening, n, hour WHERE opening.hour_start <= hour AND opening.hour_end >= hour "
	}

	query += "RETURN n.ns1__climateShelterId AS result"
	driver, err := database.NewNeo4jDriver()
	if err != nil {
		return nil, err
	}
	defer database.CloseNeo4jDriver(driver)

	session := driver.NewSession(context.Background(), neo4j.SessionConfig{})
	defer database.CloseNeo4jSession(session)

	layout := time.RFC3339Nano
	t, err := time.Parse(layout, date)

	if err != nil {
		panic(err)
	}
	currentDate := t.Format("2006-01-02")
	weekday := t.Weekday().String()
	hour := t.Format("15:04")

	params := map[string]any{
		"selectedCharacteristics": selectedCharacteristics,
		"selectedBuildingType":    selectedBuildingType,
		"weekday":                 weekday,
		"currentDate":             currentDate,
		"hour":                    hour,
	}

	result, err := session.Run(context.Background(), query, params)
	if err != nil {
		return nil, err
	}
	finalResult := []string{}
	for result.Next(context.Background()) {
		record := result.Record()

		value, ok := record.Get("result")
		if ok {
			strValue, ok := value.(string)
			if ok {
				finalResult = append(finalResult, strValue)
			}
		}

	}

	if err := result.Err(); err != nil {
		return nil, err
	}

	return finalResult, nil
}
