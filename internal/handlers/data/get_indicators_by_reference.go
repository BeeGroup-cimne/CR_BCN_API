package data

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"go-api/internal/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"time"
)

type Indicator struct {
	Reference       string                 `json:"reference"`
	CalculationDate string                 `json:"calculation_date"`
	Indicators      map[string]interface{} `json:"indicators"`
	Location        map[string]interface{} `json:"location"`
}

func GetIndicatorsByReferenceHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	reference := vars["reference"]
	aggregationLevel := vars["aggregation_level"]

	var collection *mongo.Collection
	switch aggregationLevel {
	case "building":
		collection = database.MongoClient.Database("cr_bcn_kpis").Collection("indicators_by_reference") //buildings_indicators_by_reference
	case "census_tract":
		collection = database.MongoClient.Database("cr_bcn_kpis").Collection("census_tracts_indicators_by_reference")
	default:
		http.Error(w, "Invalid aggregation level", http.StatusBadRequest)
		return
	}

	if reference == "" {
		http.Error(w, "'reference' param missing", http.StatusBadRequest)
		return
	}
	filter := bson.M{"reference": reference}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var result Indicator
	err := collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Reference not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error querying DB", http.StatusInternalServerError)
			log.Println("DB Query Error:", err)
		}
		return
	}
	result.Reference = reference

	if err := json.NewEncoder(w).Encode(result); err != nil {
		log.Println("Error encoding response:", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}
