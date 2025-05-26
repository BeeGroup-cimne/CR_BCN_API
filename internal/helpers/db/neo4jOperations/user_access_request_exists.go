package neo4jOperations

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go-api/internal/database"
	userAccessRequestStatus "go-api/internal/models/enums/UserAccesRequestStatus"
	"strings"
)

func UserAccessRequestExists(email string) (bool, error) {
	query := "MATCH (u:UserAccessRequest {email: $email, status: $status}) RETURN u UNION MATCH (u:APIUser {email: $email}) RETURN u"
	params := map[string]interface{}{
		"email":  strings.ToLower(email),
		"status": userAccessRequestStatus.Unprocessed,
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
