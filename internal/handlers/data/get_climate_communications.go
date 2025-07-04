package data

import (
	"context"
	"encoding/json"
	"go-api/internal/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"time"
)

// GetClimateCommunicationsHandler returns active climate communications grouped by date, indicator and metadata.
// @Summary      Get active climate communications
// @Description  Returns a map of active communications grouped by date and indicator, filtered by building references.
// @Tags         Climate communications
// @Accept       json
// @Produce      json
// @Param        references  body      models.ReferenceRequest  true  "List of building references to filter by"
// @Success      200         {object}  map[string]map[string][]models.CommunicationGroup  "Grouped communications"
// @Failure      400         {object}  models.Message  "Invalid input"
// @Failure      405         {object}  models.Message  "Invalid method"
// @Failure      500         {object}  models.Message  "Internal server error"
// @Security BearerAuth
// @Router       /data/climate_communications [post]
func GetClimateCommunicationsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	type RequestBody struct {
		References []string `json:"references"`
	}

	var input RequestBody
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		log.Println("Invalid JSON body:", err)
		return
	}
	if len(input.References) == 0 {
		http.Error(w, "No references provided", http.StatusBadRequest)
		return
	}

	collection := database.MongoClient.Database("cr_bcn_kpis").Collection("climate_communications")

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
	defer cancel()

	projection := bson.D{
		{"date", 1},
		{"building_reference", 1},
		{"metadata", 1},
		{"communications", 1},
	}
	filter := bson.M{
		"building_reference": bson.M{
			"$in": input.References,
		},
	}

	cursor, err := collection.Find(ctx, filter, options.Find().SetProjection(projection).SetAllowDiskUse(true))
	if err != nil {
		http.Error(w, "Error querying DB", http.StatusInternalServerError)
		log.Println("Error querying DB:", err)
		return
	}
	defer func() {
		if err := cursor.Close(ctx); err != nil {
			log.Println("Error closing cursor:", err)
		}
	}()

	finalResult := make(map[string]map[string][]map[string]interface{})

	for cursor.Next(ctx) {
		var result bson.M
		if err := cursor.Decode(&result); err != nil {
			log.Println("Error decoding result:", err)
			http.Error(w, "Error decoding result", http.StatusInternalServerError)
			return
		}

		dateRaw, ok := result["date"].(primitive.DateTime)
		if !ok {
			log.Println("Warning: 'date' is missing or not a valid DateTime in document:", result)
			continue
		}
		date := time.UnixMilli(int64(dateRaw)).UTC().Format("2006-01-02")

		buildingReference, ok := result["building_reference"].(string)
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

		if _, exists := finalResult[date]; !exists {
			finalResult[date] = make(map[string][]map[string]interface{})
		}

		for indicator, value := range communications {
			isActive, ok := value.(bool)
			if !ok || !isActive {
				continue
			}

			found := false
			for i, group := range finalResult[date][indicator] {
				existingMetadata, ok := group["metadata"].(bson.M)
				if !ok {
					continue
				}

				if existingMetadata["mean_temperature_day"] == metadata["mean_temperature_day"] &&
					existingMetadata["mean_temperature_night"] == metadata["mean_temperature_night"] {

					finalResult[date][indicator][i]["building_references"] = append(finalResult[date][indicator][i]["building_references"].([]string), buildingReference)
					found = true
					break
				}
			}

			if !found {
				newGroup := map[string]interface{}{
					"building_references": []string{buildingReference},
					"metadata":            metadata,
				}
				finalResult[date][indicator] = append(finalResult[date][indicator], newGroup)
			}
		}
	}

	if err := cursor.Err(); err != nil {
		log.Println("Cursor iteration error:", err)
		http.Error(w, "Error iterating over results", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(finalResult); err != nil {
		log.Println("Error encoding response:", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}
