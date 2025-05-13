package models

type UpdateStatusAccessRequest struct {
	Email  string `json:"email"`
	Status string `json:"status"`
}
