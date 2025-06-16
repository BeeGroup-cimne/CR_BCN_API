package neo4jOperations

import (
	"context"
	"fmt"
	"github.com/im7mortal/UTM"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go-api/internal/database"
	"log"
	"strconv"
	_ "strconv"
	"time"
)

func GetClimateShelterCoordinates(date string, characteristics map[string]bool, buildingType map[string]bool) (map[string][]interface{}, error) {

	query := "MATCH (n:ns1__ClimateSheltersBuildingSpace)-[:geosp__hasGeometry]-(p:geo__Point) "

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
			"date($currentDate) AS currentDate, $weekday AS weekday, $hour AS hour, n, p " +
			"WITH apoc.convert.fromJsonList(cleanedJson) AS scheduleJson, n, p, currentDate, weekday, hour " +
			"UNWIND scheduleJson AS outerList " +
			"UNWIND outerList AS schedule " +
			"WITH schedule, n, p, currentDate, weekday, hour " +
			"UNWIND (CASE WHEN size(schedule.periods) = 0 THEN [NULL] ELSE schedule.periods END) AS period " +
			"WITH schedule, n, p, weekday, hour, period, currentDate " +
			"WHERE period IS NULL OR ( currentDate.month >= period.month_start AND currentDate.month <= period.month_end AND currentDate.day >= period.day_start AND currentDate.day <= period.day_end) " +
			"UNWIND schedule.weekly_profile AS profile WITH profile, n, p,  weekday, hour WHERE weekday IN profile.days " +
			"UNWIND profile.opening_hours AS opening WITH opening, n, p, hour WHERE opening.hour_start <= hour AND opening.hour_end >= hour "
	}

	query += "RETURN apoc.map.fromPairs(collect([n.ns1__climateShelterId, [p.ns1__x_coordinate, p.ns1__y_coordinate]])) AS result"
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
	finalMap := make(map[string][]interface{})
	zoneNumber := 31
	zoneLetter := "T"

	for result.Next(context.Background()) {
		record := result.Record()

		mapValue, ok := record.Values[0].(map[string]interface{})
		if !ok {
			fmt.Println("El valor no es del tipo map[string][]interface{}")
			continue
		}

		for innerKey, innerVal := range mapValue {
			println(innerKey)
			if slice, ok := innerVal.([]interface{}); ok {
				if len(slice) >= 2 {
					easting, ok1 := slice[0].(string)
					northing, ok2 := slice[1].(string)

					if ok1 && ok2 {
						f, err := strconv.ParseFloat(easting, 64)
						f2, err := strconv.ParseFloat(northing, 64)
						lat, lon, err := UTM.ToLatLon(f, f2, zoneNumber, zoneLetter)
						if err != nil {
							log.Fatalf("Error al convertir: %v", err)
						} else {
							finalMap[innerKey] = []interface{}{lat, lon}
						}

					} else {
						fmt.Println("Los elementos no son strings")
					}
				}
			} else {
				fmt.Println("innerVal no es []interface{}")
			}

			if !ok {
				continue
			}
		}
	}
	if err := result.Err(); err != nil {
		return nil, err
	}

	return finalMap, nil
}
