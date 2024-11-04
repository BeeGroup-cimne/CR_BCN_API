package models

type TokenRequest struct {
	Email string    `json:"email"`
	Role  UserRoles `json:"role"`
}
