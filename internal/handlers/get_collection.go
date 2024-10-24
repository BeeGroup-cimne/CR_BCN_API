package handlers

import (
	"context"
	"encoding/json"
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
