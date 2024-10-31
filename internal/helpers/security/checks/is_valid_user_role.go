package checks

import (
	"go-api/internal/models"
	"go-api/internal/models/enums/userRoles"
)

func IsValidUserRole(role models.UserRoles) bool {
	switch role {
	case userRoles.Admin, userRoles.Unregistered, userRoles.Registered:
		return true
	default:
		return false
	}
}
