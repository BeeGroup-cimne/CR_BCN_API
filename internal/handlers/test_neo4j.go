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

	driver, err := database.NewNeo4jDriver()
	if err != nil {
		log.Fatalf("Error when connecting to db: %s", err)
	}
	defer database.CloseNeo4jDriver(driver)

	session := driver.NewSession(context.Background(), neo4j.SessionConfig{})
	defer database.CloseNeo4jSession(session)

	result, err := session.Run(context.Background(), "MATCH (n:s4bldg__Building) RETURN n", nil)
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
	_, err = w.Write(jsonResponse)
	if err != nil {
		http.Error(w, "Failed to write json body in response.", http.StatusInternalServerError)
		return
	}
}
