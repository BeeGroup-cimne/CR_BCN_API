package checks

import (
	"errors"
	"go-api/internal/helpers/db/neo4jOperations"
	"go-api/internal/helpers/security/password"
	"go-api/internal/models"
)

func UserUpdatePasswordParamsChecks(user models.UserChangePasswordRequest) error {
	//USER CHECKS
	emailUser, err := neo4jOperations.GetUser(user.Email)
	if err != nil {
		return err
	}

	// TOKEN CHECKS
	err = password.ValidateRecoveryToken(emailUser, user.Token)
	if err != nil {
		return err
	}

	//PASSWORD CHECKS
	if user.Password == "" {
		return errors.New("No password")
	}
	if !IsSecurePassword(user.Password) {
		return errors.New("Invalid password format. The password must have at least:\n- 8 characters\n- 1 Uppercase character\n- 1 Lower case character\n- 1 Number\n- 1 special character !@#$%^&*(),.?\":{}|<>_")
	}
	return nil
}
