package checks

import (
	"errors"
	"go-api/internal/helpers/db/neo4jOperations"
	"go-api/internal/models"
)

func UserAccessRequestCreateParamsChecks(user models.UserAccessRequest) error {
	//PASSWORD CHECKS
	if user.Password == "" {
		return errors.New("No password")
	}
	if !IsSecurePassword(user.Password) {
		return errors.New("Invalid password format. The password must have at least:\n- 8 characters\n- 1 Uppercase character\n- 1 Lower case character\n- 1 Number\n- 1 special character !@#$%^&*(),.?\":{}|<>_")
	}

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
