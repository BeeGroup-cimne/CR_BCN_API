package data

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"go-api/internal/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"time"
)

func GetIndicatorsByGroupHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	lang := vars["lang"]
	aggregationLevel := vars["agregation_level"]
	var collection *mongo.Collection
	switch aggregationLevel {
	case "building":
		collection = database.MongoClient.Database("cr_bcn_kpis").Collection("buildings_indicators_by_group")
	case "census_tracts":
		collection = database.MongoClient.Database("cr_bcn_kpis").Collection("census_tracts_indicators_by_group")
	default:
		http.Error(w, "Invalid aggregation level", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	projection := bson.M{
		"key":                                1,
		fmt.Sprintf("translations.%s", lang): 1,
		"kpis.key":                           1,
		fmt.Sprintf("kpis.translations.%s", lang): 1,
	}

	cursor, err := collection.Find(ctx, bson.M{}, options.Find().SetProjection(projection))
	if err != nil {
		http.Error(w, "Error querying DB", http.StatusInternalServerError)
		log.Println("DB Query Error:", err)
		return
	}
	defer cursor.Close(ctx)

	indicatorsGrouped := make(map[string]map[string]string)

	for cursor.Next(ctx) {
		var result bson.M
		if err := cursor.Decode(&result); err != nil {
			log.Println("Error decoding result:", err)
			http.Error(w, "Error decoding result", http.StatusInternalServerError)
			return
		}

		groupKey, ok := result["key"].(string)
		if !ok {
			log.Println("Invalid group key format:", result)
			continue
		}

		groupTranslations, ok := result["translations"].(bson.M)
		if !ok {
			log.Printf("Warning: No translations found for group '%s'\n", groupKey)
			continue
		}

		groupName, exists := groupTranslations[lang]
		if !exists {
			log.Printf("Warning: No translation found for language '%s' in group '%s'\n", lang, groupKey)
			continue
		}

		rawKPIs, exists := result["kpis"]
		if !exists {
			log.Printf("Warning: No KPIs found for group '%s'\n", groupKey)
			continue
		}

		kpiArray, ok := rawKPIs.(bson.A)
		if !ok {
			log.Println("Invalid KPIs format in document:", result)
			continue
		}

		kpiMap := make(map[string]string)

		for _, kpiRaw := range kpiArray {
			kpi, ok := kpiRaw.(bson.M)
			if !ok {
				log.Println("Invalid KPI format:", kpiRaw)
				continue
			}

			kpiKey, ok := kpi["key"].(string)
			if !ok {
				log.Println("Invalid KPI key format:", kpi)
				continue
			}

			kpiTranslations, ok := kpi["translations"].(bson.M)
			if !ok {
				log.Println("Invalid KPI translations format:", kpi)
				continue
			}

			kpiName, exists := kpiTranslations[lang]
			if !exists {
				log.Printf("Warning: No translation found for KPI '%s' in language '%s'\n", kpiKey, lang)
				continue
			}

			kpiMap[fmt.Sprintf("%v", kpiName)] = kpiKey
		}

		// Agregar grupo al mapa final solo si tiene KPIs
		if len(kpiMap) > 0 {
			indicatorsGrouped[fmt.Sprintf("%v", groupName)] = kpiMap
		}
	}

	if err := json.NewEncoder(w).Encode(indicatorsGrouped); err != nil {
		log.Println("Error encoding response:", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}
