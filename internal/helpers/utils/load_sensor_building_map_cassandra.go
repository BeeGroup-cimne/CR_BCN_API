package utils

import (
	"fmt"
	"strings"

	"github.com/gocql/gocql"
	"go-api/internal/models"
)

const hashPrefix = "https://climres.cat#skyspark-measurement-"

// hash format:
// https://climres.cat#skyspark-measurement-{building}-{sensorID}-{dwelling}-{property}-{aggregation}
//
// e.g. https://climres.cat#skyspark-measurement-am44-07213834-3-1-temperatura-PT15M
//
//	https://climres.cat#skyspark-measurement-am44-fcb1a032-bol250-humedad-PT15M
//
// "dwelling" is NOT always two "-"-separated tokens (e.g. "4-1", "3-1h") —
// sometimes it's a single token (e.g. "bol250"). building, sensorID,
// property and aggregation are always exactly one token each, so we split
// from both ends and treat whatever is left in the middle as the dwelling.
func parseHash(hash string) (building, sensorID, dwelling, property, aggregation string, ok bool) {
	if !strings.HasPrefix(hash, hashPrefix) {
		return "", "", "", "", "", false
	}

	parts := strings.Split(strings.TrimPrefix(hash, hashPrefix), "-")
	if len(parts) < 4 {
		return "", "", "", "", "", false
	}

	building = parts[0]
	sensorID = parts[1]
	aggregation = parts[len(parts)-1]
	property = parts[len(parts)-2]
	dwelling = strings.Join(parts[2:len(parts)-2], "-")

	if dwelling == "" {
		return "", "", "", "", "", false
	}

	return building, sensorID, dwelling, property, aggregation, true
}

// LoadSensorBuildingMapFromCassandra derives the sensor -> building/dwelling
// map directly from the distinct partitions stored in Cassandra, instead of
// relying on a manually maintained JSON file. It scans one representative
// table (e.g. the PT15M one, which has the finest granularity and therefore
// the full set of physical sensors) since the building/dwelling of a sensor
// doesn't change with the aggregation level.
//
// The session is passed in explicitly (instead of importing the database
// package) to avoid an import cycle: database.ConnectCassandraDB calls this
// function, so this package cannot import database back.
func LoadSensorBuildingMapFromCassandra(session *gocql.Session, fullyQualifiedTable string) (map[string]models.SensorInfo, error) {

	iter := session.Query(
		fmt.Sprintf(`SELECT DISTINCT hash, year FROM %s`, fullyQualifiedTable),
	).Iter()

	sensorMap := make(map[string]models.SensorInfo)

	var hash string
	var year int

	for iter.Scan(&hash, &year) {
		building, sensorID, dwelling, property, _, ok := parseHash(hash)
		if !ok {
			continue
		}

		building = strings.ToUpper(building)
		dwellingUpper := strings.ToUpper(dwelling)

		// Composite key: same shape used by GetIndoorSensorData to
		// rebuild the hash later ("{building}-{sensorID}-{dwelling}").
		key := fmt.Sprintf("%s-%s-%s", building, sensorID, dwelling)

		sensorMap[key] = models.SensorInfo{
			Building: building,
			Dwelling: fmt.Sprintf("%s %s", building, dwellingUpper),
			Type:     property,
		}
		// Composite key: same shape used by GetIndoorSensorData to
		// rebuild the hash later ("{sensorID}-{dwelling}").
		//key := fmt.Sprintf("%s-%s", sensorID, dwelling)
		//
		//sensorMap[key] = models.SensorInfo{
		//	Building: building,
		//	Dwelling: fmt.Sprintf("%s %s", building, dwellingUpper),
		//	Type:     property, // last property seen for this device; informational only
		//}
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}

	return sensorMap, nil
}
