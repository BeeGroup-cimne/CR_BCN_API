package models

type UserAccessRequest struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	Entity     string `json:"entity"`
	Department string `json:"department"`
	Name       string `json:"name"`
}
