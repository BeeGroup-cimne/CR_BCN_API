package models

type WebResource struct {
	Id       *int     `json:"id"`
	Type     string   `json:"type"`
	Value    string   `json:"value"`
	Labels   []string `json:"labels"`
	Disabled bool     `json:"disabled"`
}
