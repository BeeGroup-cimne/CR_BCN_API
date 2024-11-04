package checks

import (
	"errors"
	"fmt"
	"go-api/internal/helpers/db/neo4jOperations"
	"go-api/internal/models"
	"time"
)

func UserCreateParamsChecks(user models.CreateUserRequest) error {
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
	userExists, err := neo4jOperations.UserExists(user.Email)
	if err != nil {
		return err
	}
	if userExists {
		return errors.New("Email already in use.")
	}

	//ROLE CHECKS
	if user.Role == "" {
		return errors.New("No role")
	}
	if !IsValidUserRole(user.Role) {
		return errors.New("Invalid role.")
	}

	//DATE CHECKS
	if user.ExpirationDate != nil && !IsValidDate(*user.ExpirationDate) {
		return errors.New(fmt.Sprintf("The expiration date format must be the following: %s", time.RFC3339))
	}
	if user.StartDate != nil && !IsValidDate(*user.StartDate) {
		return errors.New(fmt.Sprintf("The start date format must be the following: %s", time.RFC3339))
	}
	if user.ExpirationDate != nil && user.StartDate != nil {
		startDate, _ := time.Parse(time.RFC3339, *user.StartDate)
		expiration, _ := time.Parse(time.RFC3339, *user.ExpirationDate)
		if expiration.Before(startDate) {
			return errors.New("Expiration date must be after the start date")
		}
	}

	//CONNECTIONS CHECKS
	if user.MaxDailyConnections != nil && *user.MaxDailyConnections < 0 {
		return errors.New("Max daily connections must be greater than or equal to 0")
	}
	if user.MaxHourlyConnections != nil && *user.MaxDailyConnections < 0 {
		return errors.New("Max hourly connections must be greater than or equal to 0")
	}
	if user.MaxHourlyConnections != nil && user.MaxDailyConnections != nil && *user.MaxDailyConnections < *user.MaxHourlyConnections {
		return errors.New("Max daily connections must be greater than or equal to max hourly connections")
	}
	return nil
}
