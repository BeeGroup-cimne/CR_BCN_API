package neo4jOperations

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go-api/internal/database"
	"go-api/internal/helpers/transform"
	"go-api/internal/helpers/utils"
	"go-api/internal/models/enums/timePrecision"
	"time"
)

func UpdateVisitCount(id int, lastConnectionDate *string) error {
	now := time.Now().UTC()
	lastConnection := utils.TernaryOperator(lastConnectionDate == nil, now, transform.ConvertIsoStringToTime(transform.ConvertNullableToString(lastConnectionDate)))
	newHour := !utils.IsEqualDate(lastConnection, now, timePrecision.Hour)
	newDay := !utils.IsEqualDate(lastConnection, now, timePrecision.Day)
	hourStr := utils.TernaryOperator[string](newHour, "1", "u.accumulated_hourly_connections + 1")
	dayStr := utils.TernaryOperator[string](newDay, "1", "u.accumulated_daily_connections + 1")

	query := fmt.Sprintf(`
MATCH (u)
WHERE id(u)=$id
SET u.last_connection_date=$last_connection,u.accumulated_hourly_connections=%s,u.accumulated_daily_connections=%s, u.accumulated_total_connections=u.accumulated_total_connections+1
`, hourStr, dayStr)

	params := map[string]interface{}{
		"id":              id,
		"last_connection": time.Now().UTC().Format(time.RFC3339),
	}

	driver, err := database.NewNeo4jDriver()
	if err != nil {
		return err
	}
	defer database.CloseNeo4jDriver(driver)

	session := driver.NewSession(context.Background(), neo4j.SessionConfig{})
	defer database.CloseNeo4jSession(session)

	result, err := session.Run(context.Background(), query, params)
	if err != nil {
		return fmt.Errorf("query execution failed: %w", err)
	}

	if !result.Next(context.Background()) {
		if result.Err() != nil {
			return fmt.Errorf("error iterating results: %w", result.Err())
		}
		return fmt.Errorf("user not found with id: %d", id)
	}

	return nil
}
