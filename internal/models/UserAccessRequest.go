package models

type UserAccessRequest struct {
	Id         int64                   `json:"id"`
	Email      string                  `json:"email"`
	Entity     string                  `json:"entity"`
	Department string                  `json:"department"`
	Name       string                  `json:"name"`
	Status     UserAccessRequestStatus `json:"status"`
	//Password   string                  `json:"password"`
}
