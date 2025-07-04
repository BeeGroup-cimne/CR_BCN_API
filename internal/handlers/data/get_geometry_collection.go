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
	"os"
	"strings"
	"time"
)

type GeoJSON struct {
	Type     string    `json:"type"`
	CRS      *CRS      `json:"crs,omitempty"`
	Features []Feature `json:"features"`
}

type CRS struct {
	Type       string            `json:"type"`
	Properties map[string]string `json:"properties"`
}

type Feature struct {
	Type       string     `json:"type"`
	Properties Properties `json:"properties"`
	Geometry   Geometry   `json:"geometry"`
}

type Properties struct {
	FID       int         `json:"fid"`
	Reference string      `json:"reference"`
	KPIs      interface{} `json:"kpis,omitempty"`
}

type Geometry struct {
	Type        string      `json:"type"`
	Coordinates interface{} `json:"coordinates"`
}

type MongoDoc struct {
	Reference string                 `bson:"reference"`
	KPIs      map[string]interface{} `bson:"kpis"`
}

const geoJSONFilePath = "internal/data/residential_buildings_bcn.geojson"

// GetGeometryCollectionHandler godoc
// @Summary Get a GeoJSON file with KPIs merged from the specified MongoDB collection
// @Description Returns a GeoJSON collection where each feature includes the KPI values from the given collection. If the collection name ends with `_data`, it returns all data. Otherwise, it filters by `calculation_date`.
// @Tags Indicators
// @Accept json
// @Produce application/geo+json
// @Param collection_name path string true "Name of the MongoDB collection"
// @Param calculation_date query string false "Calculation date to filter the KPIs (required if collection name does not end with '_data')"
// @Success 200 {object} GeoJSON
// @Failure 400 {string}  string  "No data found for the given date"
// @Failure 500 {string}  string  "Internal server error while generating GeoJSON"
// @Security BearerAuth
// @Router /data/geometry_collection/{collection_name} [get]
func GetGeometryCollectionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	calculation_date := r.URL.Query().Get("calculation_date")

	vars := mux.Vars(r)
	collection_name := vars["collection_name"]

	//read GeoJSON file
	geojson, err := loadGeoJSONFromFile(geoJSONFilePath)
	if err != nil {
		http.Error(w, "Error loading GeoJSON file", http.StatusInternalServerError)
		return
	}

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

	kpisMap := make(map[string]interface{})

	for cursor.Next(ctx) {
		var doc MongoDoc
		if err := cursor.Decode(&doc); err != nil {
			log.Fatal(err)
		}

		for ref, value := range doc.KPIs {
			kpisMap[ref] = value
		}
	}

	for i, feature := range geojson.Features {
		ref := feature.Properties.Reference

		if kpiValue, exists := kpisMap[ref]; exists {
			geojson.Features[i].Properties.KPIs = &kpiValue
		} else {
			geojson.Features[i].Properties.KPIs = nil
		}
	}

	// From GeoJSON to JSON
	geoJSONBytes, err := json.MarshalIndent(geojson, "", "  ")
	if err != nil {
		http.Error(w, "Error generating GeoJSON", http.StatusInternalServerError)
		return
	}
	// Allow downloads
	w.Header().Set("Content-Type", "application/geo+json")
	w.Header().Set("Content-Disposition", "attachment; filename=geojson_data.geojson")
	w.Write(geoJSONBytes)
}

func loadGeoJSONFromFile(filePath string) (*GeoJSON, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var geojson GeoJSON
	if err := json.Unmarshal(data, &geojson); err != nil {
		return nil, err
	}

	return &geojson, nil
}
