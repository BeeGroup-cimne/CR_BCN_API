package neo4jOperations

import (
	"context"
	"fmt"
	"github.com/im7mortal/UTM"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go-api/internal/database"
	userAccessRequestStatus "go-api/internal/models/enums/UserAccesRequestStatus"
	"log"
	"strconv"
	_ "strconv"
)

func GetClimateShelterCoordinates() (map[string][]interface{}, error) {

	query := "MATCH (n:ns1__ClimateSheltersBuildingSpace)-[:geosp__hasGeometry]-(p:geo__Point) " +
		"RETURN apoc.map.fromPairs(collect([n.ns1__climateShelterId, [p.ns1__x_coordinate, p.ns1__y_coordinate]])) " +
		"AS result"

	driver, err := database.NewNeo4jDriver()
	if err != nil {
		return nil, err
	}
	defer database.CloseNeo4jDriver(driver)

	session := driver.NewSession(context.Background(), neo4j.SessionConfig{})
	defer database.CloseNeo4jSession(session)

	params := map[string]interface{}{
		"status": userAccessRequestStatus.Unprocessed,
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
