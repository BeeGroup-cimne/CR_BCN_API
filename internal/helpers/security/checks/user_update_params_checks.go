package checks

import (
	"errors"
	"fmt"
	"go-api/internal/helpers/db/neo4jOperations"
	"go-api/internal/helpers/security/password"
	"go-api/internal/helpers/utils"
	"go-api/internal/models"
	"time"
)

func UserUpdateParamsChecks(user models.User, userUpdate models.UpdateUserRequest) error {
	//PASSWORD CHECKS
	if userUpdate.Password != nil {
		if *userUpdate.Password == "" {
			return errors.New("No password")
		}
		if password.CheckPasswordHash(*userUpdate.Password, *user.Password) {
			return errors.New("Password already in use")
		}
		if !IsSecurePassword(*userUpdate.Password) {
			return errors.New("Invalid password format. The password must have at least:\n- 8 characters\n- 1 Uppercase character\n- 1 Lower case character\n- 1 Number\n- 1 special character !@#$%^&*(),.?\":{}|<>_")
		}
	}

	//USER CHECKS
	if userUpdate.Email != nil {
		if *userUpdate.Email == "" {
			return errors.New("No email")
		}
		if !IsValidEmail(*userUpdate.Email) {
			return errors.New("Invalid email format.")
		}
		userExists, err := neo4jOperations.UserExists(*userUpdate.Email)
		if err != nil {
			return err
		}
		if userExists {
			return errors.New("Email already in use.")
		}
	}

	if userUpdate.Role != nil {
		//ROLE CHECKS
		if *userUpdate.Role == "" {
			return errors.New("No role")
		}
		if !IsValidUserRole(*userUpdate.Role) {
			return errors.New("Invalid role.")
		}
	}

	//DATE CHECKS
	if userUpdate.ExpirationDate != nil && *userUpdate.ExpirationDate != "" && !IsValidDate(*userUpdate.ExpirationDate) {
		return errors.New(fmt.Sprintf("The expiration date format must be the following: %s", time.RFC3339))
	}
	if userUpdate.StartDate != nil && *userUpdate.StartDate != "" && !IsValidDate(*userUpdate.StartDate) {
		return errors.New(fmt.Sprintf("The start date format must be the following: %s", time.RFC3339))
	}
	if userUpdate.ExpirationDate != nil && *userUpdate.ExpirationDate != "" && userUpdate.StartDate != nil && *userUpdate.StartDate != "" {
		startDate, _ := time.Parse(time.RFC3339, *userUpdate.StartDate)
		expiration, _ := time.Parse(time.RFC3339, *userUpdate.ExpirationDate)
		if expiration.Before(startDate) {
			return errors.New("Expiration date must be after the start date")
		}
	}

	//CONNECTIONS CHECKS
	if userUpdate.MaxDailyConnections != nil && *userUpdate.MaxDailyConnections < 0 {
		return errors.New("Max daily connections must be greater than or equal to 0")
	}
	if userUpdate.MaxHourlyConnections != nil && *userUpdate.MaxDailyConnections < 0 {
		return errors.New("Max hourly connections must be greater than or equal to 0")
	}
	if userUpdate.MaxHourlyConnections != nil && userUpdate.MaxDailyConnections != nil && *userUpdate.MaxDailyConnections < *userUpdate.MaxHourlyConnections {
		return errors.New("Max daily connections must be greater than or equal to max hourly connections")
	}

	if utils.AreAllFieldsNil(userUpdate) {
		return errors.New("No valid params in body")
	}
	return nil
}
