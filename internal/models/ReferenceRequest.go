package models

type ReferenceRequest struct {
	References []string `json:"references" example:"[0005401DF3800E,0005403DF3800E,0005404DF3800E]"`
}
