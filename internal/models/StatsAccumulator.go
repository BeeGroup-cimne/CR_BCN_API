package models

type StatsAccumulator struct {
	Values []float64
	Sum    float64
	Min    float64
	Max    float64
}
