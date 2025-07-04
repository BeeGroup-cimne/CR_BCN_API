package models

type Indicator struct {
	Reference string                 `json:"reference"`
	KPIs      map[string]interface{} `json:"kpis"`
}
