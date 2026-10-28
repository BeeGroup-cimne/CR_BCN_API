package utils

import (
	"go-api/internal/models"
	"sort"
	"strconv"
	"strings"
	"time"
)

func ComputeIndoorIndicatorsGrouped(
	rows []models.SensorRow,
	aggregation string,
	sensorMap map[string]models.SensorInfo,
	groupBy string,
) map[string][]models.IndoorStats {

	var duration time.Duration

	switch aggregation {
	case "PT15M":
		duration = 15 * time.Minute
	case "PT1H":
		duration = time.Hour
	case "P1D":
		duration = 24 * time.Hour
	default:
		duration = 15 * time.Minute
	}

	buckets := make(map[string]map[int64]*models.StatsAccumulator)

	for _, r := range rows {

		idx := strings.IndexByte(r.RowKey, '~')
		if idx == -1 {
			continue
		}

		ts, err := strconv.ParseInt(r.RowKey[:idx], 10, 64)
		if err != nil {
			continue
		}

		// RowKey: "{ts}~p:{building}:r:{building}-{sensorID}-{dwelling}"
		rest := r.RowKey[idx+1:]
		sepIdx := strings.Index(rest, ":r:")
		if sepIdx == -1 {
			continue
		}
		compositeKey := rest[sepIdx+len(":r:"):]

		info, ok := sensorMap[compositeKey]
		if !ok {
			continue
		}
		//idx := strings.IndexByte(r.RowKey, '~')
		//if idx == -1 {
		//	continue
		//}
		//
		//ts, err := strconv.ParseInt(r.RowKey[:idx], 10, 64)
		//if err != nil {
		//	continue
		//}
		//
		//sensorID := strings.TrimPrefix(r.RowKey[idx+1:], "p:am44:r:")
		//
		//info, ok := sensorMap[sensorID]
		//if !ok {
		//	continue
		//}

		t := time.Unix(ts, 0)
		hour := t.Hour()

		isHSensor := strings.HasSuffix(info.Dwelling, "H")

		if groupBy == "building" {

			if isHSensor && hour >= 8 && hour < 20 {
				continue
			}
			if !isHSensor && (hour < 8 || hour >= 20) {
				continue
			}
		}
		groupKey := info.Building
		if groupBy == "dwelling" {
			groupKey = info.Dwelling
		}

		bucket := t.Truncate(duration).Unix()

		if _, ok := buckets[groupKey]; !ok {
			buckets[groupKey] = make(map[int64]*models.StatsAccumulator)
		}

		acc, ok := buckets[groupKey][bucket]
		if !ok {
			acc = &models.StatsAccumulator{
				Min: r.Value,
				Max: r.Value,
			}
			buckets[groupKey][bucket] = acc
		}

		acc.Values = append(acc.Values, r.Value)
		acc.Sum += r.Value

		if r.Value < acc.Min {
			acc.Min = r.Value
		}
		if r.Value > acc.Max {
			acc.Max = r.Value
		}
	}

	results := make(map[string][]models.IndoorStats)

	for groupKey, bucketMap := range buckets {
		if groupKey != "AM44" && groupBy == "building" {
			continue
		}
		var stats []models.IndoorStats

		for bucket, acc := range bucketMap {

			sort.Float64s(acc.Values)

			mean := acc.Sum / float64(len(acc.Values))

			stat := models.IndoorStats{
				IntervalStart: bucket,
				IntervalEnd:   bucket + int64(duration.Seconds()) - 1,
				Mean:          mean,
				Min:           acc.Min,
				Max:           acc.Max,

				P2_5:  percentile(acc.Values, 0.025),
				P5:    percentile(acc.Values, 0.05),
				P10:   percentile(acc.Values, 0.10),
				P25:   percentile(acc.Values, 0.25),
				P50:   percentile(acc.Values, 0.50),
				P75:   percentile(acc.Values, 0.75),
				P90:   percentile(acc.Values, 0.90),
				P95:   percentile(acc.Values, 0.95),
				P97_5: percentile(acc.Values, 0.975),
			}

			stats = append(stats, stat)
		}

		sort.Slice(stats, func(i, j int) bool {
			return stats[i].IntervalStart < stats[j].IntervalStart
		})

		results[groupKey] = stats
	}

	return results
}
