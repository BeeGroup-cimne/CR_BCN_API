package data

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"go-api/internal/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"time"
)

// GetIndicatorsByReferenceDataHandler returns metadata for all indicators for a given aggregation level.
// @Summary Get indicators metadata by reference
// @Description Retrieves the metadata of all indicators (e.g., name, unit, description, labels) for a specific aggregation level, translated to the specified language.
// @Tags Indicators
// @Produce json
// @Param aggregation_level path string true "Aggregation level" Enums(building)
// @Param lang path string true "Language" Enums(en, es, ca)
// @Success 200 {object} map[string]interface{} "List of indicators with multilingual fields translated"
// @Failure 400 {object} models.HTTPError "Invalid or missing parameters"
// @Failure 500 {object} models.HTTPError "Database query error"
// @Security BearerAuth
// @Router /data/collection/indicators_by_reference_data/{aggregation_level}/{lang} [get]
func GetIndicatorsByReferenceDataHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	aggregationLevel := vars["aggregation_level"]
	lang := vars["lang"]
	langIndex := map[string]int{"en": 0, "es": 1, "ca": 2}[lang]

	if lang == "" || langIndex > 2 {
		http.Error(w, "Invalid or missing lang parameter (en, es, ca)", http.StatusBadRequest)
		return
	}

	var collection *mongo.Collection
	switch aggregationLevel {
	case "building":
		collection = database.MongoClient.Database("cr_bcn_kpis").Collection("indicators_by_reference_data")
	case "census_tract":
		// collection = database.MongoClient.Database("cr_bcn_kpis").Collection("census_tracts_indicators_by_reference")
	default:
		http.Error(w, "Invalid aggregation level", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var result struct {
		Indicators []map[string]interface{} `json:"indicators"` // Cambiado a un arreglo de mapas
	}

	err := collection.FindOne(ctx, bson.M{}).Decode(&result)
	if err != nil {
		log.Println("Error querying MongoDB:", err)
		http.Error(w, "Error querying MongoDB", http.StatusInternalServerError)
		return
	}

	filteredIndicators := make([]map[string]interface{}, len(result.Indicators))

	for i, kpiMap := range result.Indicators {
		filtered := make(map[string]interface{})
		for key, val := range kpiMap {
			switch key {
			case "description", "name", "unit", "labels_names":
				list, ok := val.(primitive.A)
				if ok && len(list) > langIndex {
					filtered[key] = list[langIndex]
				}
			default:
				filtered[key] = val
			}
		}
		filteredIndicators[i] = filtered
	}

	response := map[string]interface{}{
		"indicators": filteredIndicators,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println("Error encoding response:", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}
