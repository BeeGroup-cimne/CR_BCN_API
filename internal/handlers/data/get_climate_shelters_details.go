package data

import (
	"encoding/json"
	"github.com/gorilla/mux"
	_ "github.com/im7mortal/UTM"
	"go-api/internal/helpers/db/neo4jOperations"
	_ "go.mongodb.org/mongo-driver/mongo"
	_ "go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
)

func GetClimateSheltersDetails(w http.ResponseWriter, r *http.Request) {
	var payload RequestPayload
	json.NewDecoder(r.Body).Decode(&payload)

	vars := mux.Vars(r)
	id := vars["climate_shelter_id"]

	var shelters, _ = neo4jOperations.GetFilteredClimateShelterDetailsById(id)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(shelters)
}
