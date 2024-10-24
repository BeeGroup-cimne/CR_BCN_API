package handlers

import (
	"context"
	"encoding/json"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go-api/internal/database"
	"log"
	"net/http"
)

// GetNeo4jTestHandler
func GetNeo4jTestHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	database.ConnectNeo4j()
	query := "MATCH (n:s4bldg__Building) RETURN n"
	session := database.Neo4jDriver.NewSession(context.Background(), neo4j.SessionConfig{})
	defer session.Close(context.Background())
	result, err := session.Run(context.Background(), query, nil)
	if err != nil {
		log.Fatalf("Query failed: %s", err)
	}
	var results []map[string]interface{}

	// Iterate through the result set
	for result.Next(context.Background()) {
		record := result.Record()

		// Create a map to hold dynamic key-value pairs
		row := make(map[string]interface{})

		// Loop through the keys (fields) in the record
		for _, key := range record.Keys {
			value, _ := record.Get(key)
			row[key] = value // Store each key-value pair
		}

		// Append the dynamic map (row) to the results slice
		results = append(results, row)
	}

	// Check for errors during iteration
	if err := result.Err(); err != nil {
		log.Fatal("Error iterating result: ", err)
	}

	// Convert the result to JSON
	jsonResponse, err := json.Marshal(results)
	if err != nil {
		http.Error(w, "Failed to marshal result into JSON", http.StatusInternalServerError)
		return
	}

	// Set the content type and send the response
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)

	database.CloseNeo4j()
}
