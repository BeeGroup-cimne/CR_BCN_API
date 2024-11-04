package neo4jOperations

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go-api/internal/database"
	"strings"
)

func UserExists(email string) (bool, error) {
	query := "MATCH (u:APIUser) WHERE u.email=$email RETURN u"
	params := map[string]interface{}{
		"email": strings.ToLower(email),
	}

	driver, err := database.NewNeo4jDriver()
	if err != nil {
		return false, err
	}
	defer database.CloseNeo4jDriver(driver)

	session := driver.NewSession(context.Background(), neo4j.SessionConfig{})
	defer database.CloseNeo4jSession(session)

	result, err := session.Run(context.Background(), query, params)
	if err != nil {
		return false, err
	}
	return result.Next(context.Background()), nil
}
