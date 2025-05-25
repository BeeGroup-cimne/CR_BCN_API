package models

type UpdateStatusAccessRequest struct {
	//Email  string `json:"email"`
	Status UserAccessRequestStatus `json:"status"`
}
