package models

type CommunicationGroup struct {
	BuildingReferences []string       `json:"building_references" example:"[\"1234A\", \"5678B\"]"`
	Metadata           map[string]any `json:"metadata"`
}
