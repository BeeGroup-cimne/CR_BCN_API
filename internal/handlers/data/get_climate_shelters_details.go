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

// GetClimateSheltersDetails godoc
// @Summary Get details of a specific climate shelter
// @Description Returns detailed information about a specific climate shelter identified by ID
// @Tags Climate shelters
// @Accept json
// @Produce json
// @Param climate_shelter_id path string true "Climate Shelter ID"
// @Param payload body RequestPayload true "Request payload with date and filters"
// @Security BearerAuth
// @Router /data/climate_shelters_details/{climate_shelter_id} [post]
func GetClimateSheltersDetails(w http.ResponseWriter, r *http.Request) {
	var payload RequestPayload
	json.NewDecoder(r.Body).Decode(&payload)

	vars := mux.Vars(r)
	id := vars["climate_shelter_id"]

	var shelters, _ = neo4jOperations.GetFilteredClimateShelterDetailsById(id)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(shelters)
}
