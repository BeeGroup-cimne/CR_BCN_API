package database

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go-api/internal/helpers/env"
	"log"
	"time"
)

func NewNeo4jDriver() (neo4j.DriverWithContext, error) {
	keys := []string{"NEO4J_URI", "NEO4J_USERNAME", "NEO4J_PASSWORD", "NEO4J_TIMEOUT"}
	vars := env.GetVars(keys)
	timeout, err := time.ParseDuration(vars["NEO4J_TIMEOUT"])
	if err != nil {
		log.Fatal("Invalid token duration:", err)
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	uri := vars["NEO4J_URI"]
	username := vars["NEO4J_USERNAME"]
	password := vars["NEO4J_PASSWORD"]
	driver, err := neo4j.NewDriverWithContext(uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		log.Fatal("Error connecting to Neo4j: ", err)
		return nil, err
	}

	err = driver.VerifyConnectivity(ctx)
	if err != nil {
		log.Fatal("Error pinging Neo4j: ", err)
		return nil, err
	}

	return driver, nil
}

func CloseNeo4jDriver(driver neo4j.DriverWithContext) {
	if driver != nil {
		err := driver.Close(context.Background())
		if err != nil {
			log.Fatal("Error closing Neo4j connection: ", err)
		}
	}
}

func CloseNeo4jSession(session neo4j.SessionWithContext) {
	if session != nil {
		err := session.Close(context.Background())
		if err != nil {
			log.Fatal("Error closing Neo4j session: ", err)
		}
	}
}
