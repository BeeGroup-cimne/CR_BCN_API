package cassandra

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"go-api/internal/database"
	"go-api/internal/models"
)

// GetIndoorSensorData mirrors the HBase version's signature and output
// ([]models.SensorRow), but queries Cassandra instead.
//
// Cassandra partition key is (hash, year), so unlike HBase there is no
// free table scan by time range: we must know the exact "hash" per sensor
// and iterate year by year across [startTS, endTS].
//
// property is used to build the hash. The stored (raw) aggregation is
// always PT15M — coarser aggregations (PT1H, P1D) requested by callers are
// computed afterwards in Go by ComputeIndoorIndicatorsGrouped, not by
// switching Cassandra tables.
//
// The "am44" segment below is NOT the real building — every hash in this
// table is stored under the literal "am44" token regardless of which
// physical building the sensor actually belongs to (confirmed against
// Cassandra directly). The real building/dwelling for each sensor lives in
// database.SensorBuildingMap, loaded from internal/data/sensor_building_map_cassandra.json,
// keyed exactly as "{sensorID}-{dwelling}" (two parts, NOT three — do not
// try to split out a building segment from this key, there isn't one).
// hash: https://climres.cat#skyspark-measurement-am44-{sensorID}-{dwelling}-{property}-PT15M
//
// dwellings, when non-empty, restricts which dwellings are queried/included
// (case-insensitive match against the raw dwelling token, e.g. "3-1",
// "4-2h", "bol250" — same shape as it appears in the Cassandra hash). An
// empty/nil slice means "no filter, include every dwelling".
const storageAggregation = "PT15M"
const table = "skyspark.climres_static_pt15m"

// SensorBuildingMap keys have the shape "{sensorID}-{dwelling}", and
// sensorID never contains a "-", so splitting on the first "-" always
// isolates the dwelling portion.
func dwellingFromCompositeKey(key string) string {
	idx := strings.Index(key, "-")
	if idx == -1 {
		return ""
	}
	return key[idx+1:]
}

func GetIndoorSensorData(property string, startTS int64, endTS int64, dwellings []string) ([]models.SensorRow, error) {

	result := make([]models.SensorRow, 0, 10000)

	startYear := time.Unix(startTS, 0).UTC().Year()
	endYear := time.Unix(endTS, 0).UTC().Year()

	var dwellingFilter map[string]bool
	if len(dwellings) > 0 {
		dwellingFilter = make(map[string]bool, len(dwellings))
		for _, d := range dwellings {
			dwellingFilter[strings.ToLower(strings.TrimSpace(d))] = true
		}
	}

	for compositeKey := range database.SensorBuildingMap {

		if dwellingFilter != nil {
			dwelling := strings.ToLower(dwellingFromCompositeKey(compositeKey))
			if !dwellingFilter[dwelling] {
				continue
			}
		}

		hash := fmt.Sprintf(
			"https://climres.cat#skyspark-measurement-am44-%s-%s-%s",
			compositeKey, property, storageAggregation,
		)

		for year := startYear; year <= endYear; year++ {

			iter := database.CassandraSession.Query(
				fmt.Sprintf(`SELECT ts, info FROM %s WHERE hash = ? AND year = ? AND ts >= ? AND ts <= ?`, table),
				hash, year, startTS, endTS,
			).Iter()

			var ts int64
			var info map[string]string

			for iter.Scan(&ts, &info) {
				raw, ok := info["value"]
				if !ok {
					continue
				}
				v, err := strconv.ParseFloat(raw, 64)
				if err != nil {
					continue
				}

				// RowKey kept in the same "{ts}~p:am44:r:{sensorID}" shape
				// the aggregation code (ComputeIndoorIndicatorsGrouped) expects,
				// so downstream parsing doesn't need to change at all.
				result = append(result, models.SensorRow{
					RowKey: fmt.Sprintf("%d~p:am44:r:%s", ts, compositeKey),
					Value:  v,
				})
			}

			if err := iter.Close(); err != nil {
				return nil, err
			}
		}
	}

	return result, nil
}
