package neo4jOperations

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/im7mortal/UTM"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go-api/internal/database"
)

func GetClimateShelterCoordinates(date string, characteristics map[string]bool, buildingType map[string]bool) (map[string][]interface{}, error) {

	selectedCharacteristics := selectedKeys(characteristics)
	selectedBuildingType := selectedKeys(buildingType)

	query := "MATCH (n:ns1__ClimateSheltersBuildingSpace)-[:geosp__hasGeometry]-(p:geo__Point) " +
		climateShelterFilterConditions(selectedCharacteristics, selectedBuildingType) +
		"RETURN n.ns1__climateShelterId AS id, " +
		"p.ns1__x_coordinate AS x, p.ns1__y_coordinate AS y, " +
		"n.ns1__schedule AS schedule"

	var t time.Time
	if date != "" {
		var err error
		t, err = time.Parse(time.RFC3339Nano, date)
		if err != nil {
			return nil, err
		}
	}

	driver, err := database.NewNeo4jDriver()
	if err != nil {
		return nil, err
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
		return nil, err
	}

	finalMap := make(map[string][]interface{})
	zoneNumber := 31
	zoneLetter := "T"

	for result.Next(context.Background()) {
		record := result.Record()

		idVal, _ := record.Get("id")
		id, ok := idVal.(string)
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

		xVal, _ := record.Get("x")
		yVal, _ := record.Get("y")
		easting, ok1 := xVal.(string)
		northing, ok2 := yVal.(string)
		if !ok1 || !ok2 {
			log.Printf("refugio %s: coordenadas no son texto", id)
			continue
		}

		f, err1 := strconv.ParseFloat(easting, 64)
		f2, err2 := strconv.ParseFloat(northing, 64)
		if err1 != nil || err2 != nil {
			log.Printf("refugio %s: coordenadas no válidas (%q, %q)", id, easting, northing)
			continue
		}
		lat, lon, err := UTM.ToLatLon(f, f2, zoneNumber, zoneLetter)
		if err != nil {
			// Antes era log.Fatalf, que paraba toda la API por un solo refugio.
			log.Printf("refugio %s: error al convertir coordenadas: %v", id, err)
			continue
		}
		finalMap[id] = []interface{}{lat, lon}
	}
	if err := result.Err(); err != nil {
		return nil, err
	}

	return finalMap, nil
}
