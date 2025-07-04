package models

type HTTPError struct {
	Message string `json:"message" example:"internal server error"`
}
