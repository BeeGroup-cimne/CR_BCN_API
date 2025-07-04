package models

type LoginRequest struct {
	Email    string `json:"email" example:"info@cimne.upc.edu"`
	Password string `json:"password" example:"1234.!AB"`
}
