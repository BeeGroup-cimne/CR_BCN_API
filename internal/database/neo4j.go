package database

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go-api/internal/helpers/env"
	"log"
	"time"
)

func newNeo4jDriver() (neo4j.DriverWithContext, error) {
	keys := []string{"NEO4J_URI", "NEO4J_USERNAME", "NEO4J_PASSWORD", "NEO4J_TIMEOUT_SECONDS"}
	vars := env.GetVars(keys)
	timeout, err := time.ParseDuration(vars["NEO4J_TIMEOUT_SECONDS"])
	if err != nil {
		log.Fatal("Invalid token duration:", err)
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
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

	log.Println("Connected to Neo4j!")
	return driver, nil
}

func closeNeo4jDriver(driver neo4j.DriverWithContext) {
	if driver != nil {
		err := driver.Close(context.Background())
		if err != nil {
			log.Fatal("Error closing Neo4j connection: ", err)
		}
	}
}

func closeNeo4jSession(session neo4j.SessionWithContext) {
	if session != nil {
		err := session.Close(context.Background())
		if err != nil {
			log.Fatal("Error closing Neo4j session: ", err)
		}
	}
}

func RunNeo4jQuery(query string, params map[string]interface{}) (neo4j.ResultWithContext, error) {
	driver, err := newNeo4jDriver()
	if err != nil {
		return nil, err
	}
	defer closeNeo4jDriver(driver)

	session := driver.NewSession(context.Background(), neo4j.SessionConfig{})
	defer closeNeo4jSession(session)

	result, err := session.Run(context.Background(), query, params)
	if err != nil {
		return nil, err
	}
	return result, nil
}
