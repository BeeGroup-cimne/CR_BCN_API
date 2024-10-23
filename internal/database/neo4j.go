package database

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go-api/internal/Environment"
	"log"
	"time"
)

var Neo4jDriver neo4j.DriverWithContext

func ConnectNeo4j() {
	keys := []string{"NEO4J_URI", "NEO4J_USERNAME", "NEO4J_PASSWORD"}
	vars := Environment.GetVars(keys)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	uri := vars["NEO4J_URI"]
	username := vars["NEO4J_USERNAME"]
	password := vars["NEO4J_PASSWORD"]
	driver, err := neo4j.NewDriverWithContext(uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		log.Fatal("Error connecting to Neo4j: ", err)
	}

	err = driver.VerifyConnectivity(ctx)
	if err != nil {
		log.Fatal("Error pinging Neo4j: ", err)
	}

	Neo4jDriver = driver
	log.Println("Connected to Neo4j!")
}

func CloseNeo4j() {
	if Neo4jDriver != nil {
		err := Neo4jDriver.Close(context.Background())
		if err != nil {
			log.Fatal("Error closing Neo4j connection: ", err)
		}
		log.Println("Disconnected from Neo4j!")
	}
}

func RunNeo4jQuery(query string) neo4j.ResultWithContext {
	if Neo4jDriver == nil {
		ConnectNeo4j()
	}
	session := Neo4jDriver.NewSession(context.Background(), neo4j.SessionConfig{})
	defer session.Close(context.Background())
	result, err := session.Run(context.Background(), query, nil)
	if err != nil {
		log.Fatalf("Query failed: %s", err)
	}
	if result.Next(context.Background()) {
		return result
	}

	if err := result.Err(); err != nil {
		log.Fatalf("Result error: %s", err)
	}

	return nil
}
