package handlers

import (
	"context"
	"encoding/json"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go-api/internal/database"
	"go-api/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"time"
)

// getCollection
func GetCollectionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	collection := database.MongoClient.Database("r1db").Collection("grossFloorArea")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// query
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, "Error al realizar la consulta a MongoDB", http.StatusInternalServerError)
		return
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {

		}
	}(cursor, ctx)

	var results []bson.M

	for cursor.Next(ctx) {
		var result bson.M
		if err := cursor.Decode(&result); err != nil {
			log.Println(err)
			err := json.NewEncoder(w).Encode(models.Message{Status: "error", Message: "Error decoding result"})
			if err != nil {
				return
			}
			return
		}
		results = append(results, result)
	}

	if err := cursor.Err(); err != nil {
		log.Println(err)
		err := json.NewEncoder(w).Encode(models.Message{Status: "error", Message: "Cursor error"})
		if err != nil {
			return
		}
		return
	}
	// Configurar el header para que el cliente sepa que la respuesta es JSON
	w.Header().Set("Content-Type", "application/json")

	// Codificamos los resultados como JSON y los enviamos
	if err := json.NewEncoder(w).Encode(results); err != nil {
		http.Error(w, "Error al codificar la respuesta en JSON", http.StatusInternalServerError)
	}
}

// GetNeo4jTestHandler
func GetNeo4jTestHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	database.ConnectNeo4j()
	query := "MATCH (n:s4bldg__Building) RETURN n LIMIT 100"
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
