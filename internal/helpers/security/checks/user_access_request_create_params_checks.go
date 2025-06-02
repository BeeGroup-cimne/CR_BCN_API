package checks

import (
	"errors"
	"go-api/internal/helpers/db/neo4jOperations"
	"go-api/internal/models"
)

func UserAccessRequestCreateParamsChecks(user models.UserAccessRequest) error {

	//USER CHECKS
	if user.Email == "" {
		return errors.New("No email")
	}
	if !IsValidEmail(user.Email) {
		return errors.New("Invalid email format.")
	}
	userExists, err := neo4jOperations.UserAccessRequestExists(user.Email)
	if err != nil {
		return err
	}
	if userExists {
		return errors.New("Email already in use.")
	}

	return nil
}
