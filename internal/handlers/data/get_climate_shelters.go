package data

import (
	"context"
	"encoding/json"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go-api/internal/database"
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
	Buildings        map[string]float64    `json:"buildings"`
	Shelters         map[string][2]float64 `json:"shelters"`
	BuildingShelters map[string]string     `json:"building_shelters"`
	NoBuilding       map[string]int        `json:"n_building_wo_shelter"`
	NumberBuildings  int                   `json:"number_buildings"`
}

type Neo4jService struct {
	driver neo4j.DriverWithContext
}

func GetClimateShelters(w http.ResponseWriter, r *http.Request) {
	var payload RequestPayload
	json.NewDecoder(r.Body).Decode(&payload)

	ctx := context.TODO()

	collection := database.MongoClient.Database("cr_bcn_kpis").Collection("climate_shelters")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, _ := collection.Find(ctx, bson.M{})
	defer cursor.Close(ctx)

	buildings := make(map[string]float64)
	shelters := make(map[string][2]float64)
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

	//neo4jService, err := NewNeo4jService("neo4j://localhost:7687", "neo4j", "password")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer neo4jService.Close()
	//
	//// Pasar el servicio al handler con closure
	//http.HandleFunc("/climate-shelter-coordinates", func(w http.ResponseWriter, r *http.Request) {
	//	ctx := r.Context()
	//	results, err := GetShelterCoordinates(ctx)
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

//func (n *Neo4jService) GetShelterCoordinates(ctx context.Context) (map[string][]float64, error) {
//	cypher := `
//		MATCH (n:ns1__ClimateSheltersBuildingSpace)-[:geosp__hasGeometry]-(p:geo__Point)
//		RETURN apoc.map.fromPairs(collect([n.ns1__climateShelterId, [p.ns1__x_coordinate, p.ns1__y_coordinate]])) AS result
//	`
//
//	session := n.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
//	defer session.Close(ctx)
//
//	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
//		record, err := tx.Run(ctx, cypher, nil)
//		if err != nil {
//			return nil, err
//		}
//
//		if record.Next(ctx) {
//			value, ok := record.Record().Get("result")
//			if !ok {
//				return nil, fmt.Errorf("missing result key")
//			}
//
//			rawMap, ok := value.(map[string]any)
//			if !ok {
//				return nil, fmt.Errorf("result is not a map")
//			}
//
//			output := make(map[string][]float64)
//			for k, v := range rawMap {
//				list, ok := v.([]interface{})
//				if !ok || len(list) != 2 {
//					continue
//				}
//				x, _ := list[0].(float64)
//				y, _ := list[1].(float64)
//				output[k] = []float64{x, y}
//			}
//			return output, nil
//		}
//
//		if err := record.Err(); err != nil {
//			return nil, err
//		}
//
//		return nil, fmt.Errorf("no result from query")
//	})
//
//	if err != nil {
//		return nil, err
//	}
//
//	return result.(map[string][]float64), nil
//}
