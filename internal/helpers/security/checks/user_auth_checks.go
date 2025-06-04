package checks

import (
	"errors"
	"fmt"
	"go-api/internal/helpers/security/password"
	"go-api/internal/helpers/transform"
	"go-api/internal/helpers/utils"
	"go-api/internal/models"
	"go-api/internal/models/enums/timePrecision"
	"strconv"
	"time"
)

func UserAuthChecks(user *models.User, pwd *string) error {
	if pwd != nil {
		if !password.CheckPasswordHash(*pwd, *user.Password) {
			return errors.New("invalid username or password")
		}
	}

	if user.Blocked != nil && *user.Blocked {
		return errors.New("user is blocked")
	}

	if user.ExpirationDate != nil && *user.ExpirationDate != "" && time.Now().UTC().After(transform.ConvertIsoStringToTime(*user.ExpirationDate)) {
		return errors.New("user no longer has access")
	}

	if user.StartDate != nil && *user.StartDate != "" && time.Now().UTC().Before(transform.ConvertIsoStringToTime(*user.StartDate)) {
		return errors.New(fmt.Sprintf("user does not have access yet. It will start on %s", *user.StartDate))
	}

	now := time.Now().UTC()
	lastConnection := utils.TernaryOperator(user.LastConnectionDate == nil, now, transform.ConvertIsoStringToTime(transform.ConvertNullableToString(user.LastConnectionDate)))
	isNewHour := utils.IsEqualDate(lastConnection, now, timePrecision.Hour)
	isNewDay := utils.IsEqualDate(lastConnection, now, timePrecision.Day)
	if user.MaxHourlyConnections != nil && *user.MaxHourlyConnections != 0 && user.AccumulatedHourlyConnections != nil && !isNewHour && *user.MaxHourlyConnections <= *user.AccumulatedHourlyConnections {
		return errors.New(fmt.Sprintf("access hourly limit of %s reached.", strconv.Itoa(*user.MaxHourlyConnections)))
	}
	if user.MaxDailyConnections != nil && *user.MaxDailyConnections != 0 && user.AccumulatedDailyConnections != nil && !isNewDay && *user.MaxDailyConnections <= *user.AccumulatedDailyConnections {
		return errors.New(fmt.Sprintf("access daily limit of %s reached.", strconv.Itoa(*user.MaxDailyConnections)))
	}

	return nil
}
