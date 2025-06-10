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
	"math"
	"net/http"
	"time"
)

type RequestPayload struct {
	Date            string   `json:"date"`
	Characteristics []string `json:"characteristics"`
	BuildingType    string   `json:"building_type"`
}

type ResponsePayload struct {
	Buildings        map[string]float64       `json:"buildings"`
	Shelters         map[string][]interface{} `json:"shelters"`
	BuildingShelters map[string]string        `json:"building_shelters"`
	NoBuilding       map[string]int           `json:"n_building_wo_shelter"`
	NumberBuildings  int                      `json:"number_buildings"`
}

type Neo4jService struct {
	driver neo4j.DriverWithContext
}

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
			minTime := math.MaxFloat64
			nearestShelter := ""

			for shelterID, raw := range timesMap {
				seconds := raw.(float64)
				minutes := seconds / 60.0
				if minutes < minTime {
					minTime = minutes
					nearestShelter = shelterID
				}
				shelterSet[shelterID] = struct{}{}
			}

			buildings[buildingID] = minTime
			buildingShelters[buildingID] = nearestShelter

			if minTime > 10 {
				noBuilding["10"]++
			}
			if minTime > 20 {
				noBuilding["20"]++
			}
		}
	}
	shelters, err := neo4jOperations.GetClimateShelterCoordinates()
	if err != nil {
		log.Fatalf("Error obteniendo shelters: %v", err)
	}
	// Pasar el servicio al handler con closure
	//http.HandleFunc("/climate-shelter-coordinates", func(w http.ResponseWriter, r *http.Request) {
	//	ctx := r.Context()
	//	results, err := service.GetShelterCoordinates(ctx)
	//	if err != nil {
	//		http.Error(w, fmt.Sprintf("Neo4j query error: %v", err), http.StatusInternalServerError)
	//		return
	//	}
	//	w.Header().Set("Content-Type", "application/json")
	//	json.NewEncoder(w).Encode(results)
	//})
	////Obtener coordenadas desde Neo4j
	//results, err := neo4jService.GetShelterCoordinates(ctx)
	//if err != nil {
	//	http.Error(w, fmt.Sprintf("Neo4j query error: %v", err), http.StatusInternalServerError)
	//	return
	//}

	res := ResponsePayload{
		Buildings:        buildings,
		Shelters:         shelters,
		BuildingShelters: buildingShelters,
		NoBuilding:       noBuilding,
		NumberBuildings:  totalBuildings,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
