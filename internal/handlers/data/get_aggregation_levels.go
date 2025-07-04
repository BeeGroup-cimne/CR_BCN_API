package data

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"go-api/internal/database"
	"go-api/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"time"
)

// GetAggregationLevelsHandler godoc
// @Summary Get available aggregation levels with translation
// @Description Returns a map of translated aggregation level names (in the specified language) mapped to their keys.
// @Tags Aggregation
// @Accept json
// @Produce json
// @Param lang path string true "Language code (e.g., 'en', 'es', 'ca')"
// @Success 200 {object} map[string]string
// @Failure 405 {object} models.HTTPError "Invalid request method"
// @Failure 500 {object} models.HTTPError "Database query or decoding error"
// @Security BearerAuth
// @Router /data/collection/aggregation_levels/{lang} [get]
func GetAggregationLevelsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	lang := vars["lang"]

	collection := database.MongoClient.Database("cr_bcn_kpis").Collection("aggregation_levels")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	projection := bson.D{
		{"key", 1},
		{fmt.Sprintf("translations.%s", lang), 1},
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

	indicatorsMap := make(map[string]string)

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

		key := result["key"]
		translations, ok := result["translations"].(map[string]interface{})
		if !ok {
			translationsGeneric := result["translations"].(bson.M) // bson.M funciona como un map[string]interface{}
			translations = translationsGeneric
		}

		text, exists := translations[lang]
		if !exists {
			log.Printf("Warning: No translation found for language '%s' in document: %v\n", lang, result)
			continue
		}

		indicatorsMap[fmt.Sprintf("%v", text)] = fmt.Sprintf("%v", key)
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(indicatorsMap); err != nil {
		log.Println("Error encoding response:", err)
		return
	}
}
