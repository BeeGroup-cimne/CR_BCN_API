package data

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"go-api/internal/database"
	"go-api/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"strings"
	"time"
)

// getCollection
func GetCollectionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	calculation_date := r.URL.Query().Get("calculation_date")

	vars := mux.Vars(r)
	collection_name := vars["collection_name"]

	collection := database.MongoClient.Database("cr_bcn_kpis").Collection(collection_name)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var filter bson.M
	// query
	if strings.HasSuffix(collection_name, "_data") {
		filter = bson.M{}
	} else {
		if calculation_date == "" {
			http.Error(w, "'calculation_date' param missing", http.StatusBadRequest)
			return
		} else {
			filter = bson.M{"calculation_date": calculation_date}
		}

	}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		http.Error(w, "Error when querying DB", http.StatusInternalServerError)
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
		http.Error(w, "Error encoding response in JSON", http.StatusInternalServerError)
	}
}
