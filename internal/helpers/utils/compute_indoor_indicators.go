package utils

import (
	"go-api/internal/models"
	"sort"
	"strconv"
	"strings"
	"time"
)

func ComputeIndoorIndicators(rows []models.SensorRow, aggregation string) []models.IndoorStats {

	var duration time.Duration

	switch aggregation {
	case "15m":
		duration = 15 * time.Minute
	case "1h":
		duration = time.Hour
	case "1d":
		duration = 24 * time.Hour
	default:
		duration = 15 * time.Minute
	}

	buckets := map[int64]*models.StatsAccumulator{}

	for _, r := range rows {

		parts := strings.Split(r.RowKey, "~")
		if len(parts) != 2 {
			continue
		}

		ts, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			continue
		}

		t := time.Unix(ts, 0)

		bucket := t.Truncate(duration).Unix()

		acc, ok := buckets[bucket]

		if !ok {
			acc = &models.StatsAccumulator{
				Min: r.Value,
				Max: r.Value,
			}
			buckets[bucket] = acc
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

	var results []models.IndoorStats

	for bucket, acc := range buckets {

		sort.Float64s(acc.Values)

		mean := acc.Sum / float64(len(acc.Values))

		stats := models.IndoorStats{
			IntervalStart: bucket,
			IntervalEnd:   bucket + int64(duration.Seconds()),

			Mean: mean,
			Min:  acc.Min,
			Max:  acc.Max,

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

		results = append(results, stats)
	}

	return results
}

func percentile(values []float64, p float64) float64 {

	if len(values) == 0 {
		return 0
	}

	pos := p * float64(len(values)-1)
	i := int(pos)
	frac := pos - float64(i)

	if i+1 < len(values) {
		return values[i]*(1-frac) + values[i+1]*frac
	}

	return values[i]
}
