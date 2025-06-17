package models

type WebResourceUpdate struct {
	Type     *string   `json:"type"`
	Value    *string   `json:"value"`
	Labels   *[]string `json:"labels"`
	Disabled *bool     `json:"disabled"`
}
