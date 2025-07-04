package data

import (
	"context"
	"encoding/json"
	_ "github.com/im7mortal/UTM"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go-api/internal/database"
	"go-api/internal/helpers/db/neo4jOperations"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	_ "go.mongodb.org/mongo-driver/mongo"
	_ "go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"time"
)

type RequestPayload struct {
	Date            string          `json:"date" example:"2024-08-15"`
	Characteristics map[string]bool `json:"characteristics" swaggertype:"object" example:"{\"medical_services\":false,\"wc_wheelchair\":false,\"nurse_services\":false}"`
	BuildingType    map[string]bool `json:"building_type" swaggertype:"object" example:"{\"Senior Center\":false,\"Museum\":true}"`
}

type ResponsePayload struct {
	Buildings           map[string]float64       `json:"buildings"`
	Shelters            map[string][]interface{} `json:"shelters"`
	BuildingShelters    map[string]string        `json:"building_shelters"`
	NoBuilding          map[string]int           `json:"n_building_wo_shelter"`
	NumberBuildings     int                      `json:"number_buildings"`
	NumberShelters      int                      `json:"number_shelters"`
	NumberTotalShelters int                      `json:"number_total_shelters"`
}

type Neo4jService struct {
	driver neo4j.DriverWithContext
}

// GetClimateShelters returns accessibility data to climate shelters for each building.
//
// @Summary      Get accessibility to climate shelters
// @Description  Returns the minimum walking time (in minutes) from each building to the nearest filtered climate shelter, as well as metadata including which shelter is closest and how many buildings are beyond a 10 and 20 minute threshold.
// @Tags         Climate shelters
// @Accept       json
// @Produce      json
// @Param        payload  body      RequestPayload  true  "Filter parameters for shelter selection"
// @Success      200  {object}  ResponsePayload
// @Failure      400  {string}  string  "Invalid request payload"
// @Failure      500  {string}  string  "Internal server error"
// @Security BearerAuth
// @Router       /data/climate_shelters [post]
func GetClimateShelters(w http.ResponseWriter, r *http.Request) {
	var payload RequestPayload
	json.NewDecoder(r.Body).Decode(&payload)

	ctx := context.TODO()

	collection := database.MongoClient.Database("cr_bcn_kpis").Collection("climate_shelters")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cursor, _ := collection.Find(ctx, bson.M{})
	defer cursor.Close(ctx)

	buildings := make(map[string]float64)
	//shelters := make(map[string][2]float64)
	buildingShelters := make(map[string]string)
	noBuilding := map[string]int{"10": 0, "20": 0}
	totalBuildings := 0

	shelterSet := map[string]struct{}{}
	filtered_shelters, filtered_n_shelters, total_n_shelters, err := neo4jOperations.GetFilteredClimateShelter(payload.Date, payload.Characteristics, payload.BuildingType)
	print(filtered_shelters)
	filteredSheltersSet := make(map[string]struct{}, len(filtered_shelters))
	for _, shelter := range filtered_shelters {
		filteredSheltersSet[shelter] = struct{}{}
	}
	for cursor.Next(ctx) {
		var doc bson.M
		cursor.Decode(&doc)

		totalBuildings++
		for buildingID, v := range doc {
			if buildingID == "_id" {
				continue
			}
			timesMap, ok := v.(primitive.M)
			if !ok {
				log.Printf("El valor de %s no es del tipo esperado primitive.M", buildingID)
				continue
			}
			minTime := 40.0
			//minTime := math.MaxFloat64
			nearestShelter := ""

			for shelterID, raw := range timesMap {
				seconds := raw.(float64)
				minutes := seconds / 60.0
				//if buildingID == "0004701DF3800C" {
				//	print(minutes)
				//}
				if _, exists := filteredSheltersSet[shelterID]; exists && minutes < minTime {
					minTime = minutes
					nearestShelter = shelterID
				}
				shelterSet[shelterID] = struct{}{}
			}

			buildings[buildingID] = minTime
			buildingShelters[buildingID] = nearestShelter

			if minTime >= 10 {
				noBuilding["10"]++
			}
			if minTime >= 20 {
				noBuilding["20"]++
			}
		}
	}
	shelters, err := neo4jOperations.GetClimateShelterCoordinates(payload.Date, payload.Characteristics, payload.BuildingType)
	if err != nil {
		log.Fatalf("Error obteniendo shelters: %v", err)
	}

	res := ResponsePayload{
		Buildings:           buildings,
		Shelters:            shelters,
		BuildingShelters:    buildingShelters,
		NoBuilding:          noBuilding,
		NumberBuildings:     totalBuildings,
		NumberShelters:      filtered_n_shelters,
		NumberTotalShelters: total_n_shelters,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
