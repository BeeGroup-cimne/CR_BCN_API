package dbOperations

import (
	"fmt"
	"go-api/internal/database"
	"go-api/internal/helpers/utils"
	"go-api/internal/models"
	"time"
)

func UpdateVisitCount(usr models.User, tkn *models.TokensResponse, newHour bool, newDay bool) error {
	hourStr := utils.TernaryOperator[string](newHour, "1", "u.accumulated_hourly_connections + 1")
	dayStr := utils.TernaryOperator[string](newDay, "1", "u.accumulated_daily_connections + 1")
	tokenStr := ""

	params := map[string]interface{}{
		"id":              usr.Id,
		"last_connection": time.Now().UTC().Format(time.RFC3339),
	}

	if tkn != nil {
		params["token"] = tkn.Token
		params["refresh_token"] = tkn.RefreshToken
		tokenStr = "u.token=$token, u.refresh_token=$refresh_token"
	}

	query := fmt.Sprintf(`
MATCH (u)
WHERE id(u)=$id
SET u.last_connection_date=$last_connection,u.accumulated_hourly_connections=%s,u.accumulated_daily_connections=%d,%f
`, hourStr, dayStr, tokenStr)

	_, err := database.RunNeo4jQuery(query, params)

	if err != nil {
		return err
	}
	return nil
}
