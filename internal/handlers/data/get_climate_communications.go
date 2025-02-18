package data

import (
	"context"
	"encoding/json"
	"go-api/internal/database"
	"go-api/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"time"
)

func GetClimateCommunicationsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	collection := database.MongoClient.Database("cr_bcn_kpis").Collection("climate_communications")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	projection := bson.D{
		{"date", 1},
		{"building_reference", 1},
		{"metadata", 1},
		{"communications", 1},
	}
	filter := bson.D{}

	cursor, err := collection.Find(ctx, filter, options.Find().SetProjection(projection))
	if err != nil {
		http.Error(w, "Error when querying DB", http.StatusInternalServerError)
		return
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {

		}
	}(cursor, ctx)

	indicatorsMap := make(map[string]map[string]map[string]interface{})

	for cursor.Next(ctx) {
		var result bson.M
		if err := cursor.Decode(&result); err != nil {
			log.Println("Error decoding result:", err)
			err := json.NewEncoder(w).Encode(models.Message{Status: "error", Message: "Error decoding result"})
			if err != nil {
				return
			}
			return
		}

		dateRaw, ok := result["date"].(primitive.DateTime)
		if !ok {
			log.Println("Warning: 'date' is missing or not a valid DateTime in document:", result)
			continue
		}
		date := time.UnixMilli(int64(dateRaw)).UTC().Format("2006-01-02") // Convertir a formato YYYY-MM-DD

		building_reference, ok := result["building_reference"].(string)
		if !ok {
			log.Println("Warning: 'building_reference' is missing or not a string in document:", result)
			continue
		}

		communications, ok := result["communications"].(bson.M)
		if !ok {
			log.Println("Warning: 'communications' is missing or not a valid map in document:", result)
			continue
		}

		metadata, ok := result["metadata"].(bson.M)
		if !ok {
			log.Println("Warning: 'metadata' is missing or not a valid map in document:", result)
			continue
		}

		if _, exists := indicatorsMap[date]; !exists {
			indicatorsMap[date] = make(map[string]map[string]interface{})
		}
		indicatorsMap[date][building_reference] = make(map[string]interface{})
		indicatorsMap[date][building_reference]["communications"] = communications
		indicatorsMap[date][building_reference]["metadata"] = metadata
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(indicatorsMap); err != nil {
		log.Println("Error encoding response:", err)
		return
	}
}
