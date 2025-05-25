package userAccessRequestStatus

import "go-api/internal/models"

const (
	Unprocessed models.UserAccessRequestStatus = "unprocessed"
	Approved    models.UserAccessRequestStatus = "approved"
	Declined    models.UserAccessRequestStatus = "declined"
)

var AllStatus = []models.UserAccessRequestStatus{
	Unprocessed,
	Approved,
	Declined,
}
