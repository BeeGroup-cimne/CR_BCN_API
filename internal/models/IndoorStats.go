package models

type IndoorStats struct {
	IntervalStart int64  `json:"interval_start"`
	IntervalEnd   int64  `json:"interval_end"`
	Aggregation   string `json:"aggregation"`

	Mean float64 `json:"mean"`
	Min  float64 `json:"min"`
	Max  float64 `json:"max"`

	P2_5  float64 `json:"p2_5"`
	P5    float64 `json:"p5"`
	P10   float64 `json:"p10"`
	P25   float64 `json:"p25"`
	P50   float64 `json:"p50"`
	P75   float64 `json:"p75"`
	P90   float64 `json:"p90"`
	P95   float64 `json:"p95"`
	P97_5 float64 `json:"p97_5"`
}
