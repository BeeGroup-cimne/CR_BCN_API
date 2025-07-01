package data

import (
	"context"
	"encoding/json"
	_ "github.com/im7mortal/UTM"
	"go-api/internal/database"
	"go.mongodb.org/mongo-driver/bson"
	_ "go.mongodb.org/mongo-driver/mongo"
	_ "go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"time"
)

func GetClimateSheltersInfo(w http.ResponseWriter, r *http.Request) {
	var payload RequestPayload
	json.NewDecoder(r.Body).Decode(&payload)

	ctx := context.TODO()
	collection := database.MongoClient.Database("cr_bcn_kpis").Collection("climate_shelters")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cursor, _ := collection.Find(ctx, bson.M{})
	defer cursor.Close(ctx)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("")
}
