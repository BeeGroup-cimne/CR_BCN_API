package data

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"go-api/internal/database"
	"go-api/internal/helpers/db/hbase"
	"go-api/internal/helpers/utils"
	"net/http"
	"strconv"
)

// GetIndoorIndicators returns indoor sensor indicators aggregated by time and level.
// @Summary      Get indoor sensor indicators
// @Description  Returns indoor sensor indicators for a property within a time range.
// @Tags         Indicators
// @Param        property           path      string  true   "Property name"
// @Param        start              path      integer true   "Start timestamp (unix)"
// @Param        end                path      integer true   "End timestamp (unix)"
// @Param        aggregation_level  query     string  false  "Aggregation level (dwelling or building). Default: building"
// @Param        aggregation_time   query     string  false  "Temporal aggregation (e.g., PT15M, PT1H, P1D)"
// @Success      200  {array}       models.IndoorStats             "Aggregated indoor indicators"
// @Failure      400  {object}      models.Message     "Missing or invalid parameters"
// @Failure      405  {object}      models.Message     "Invalid HTTP method"
// @Failure      500  {object}      models.Message     "Internal server error"
// @Security     BearerAuth
// @Router       /data/indoor_sensor_data/{property}/{aggregation_level}/{aggregation_time}/{start}:{end} [get]
func GetIndoorIndicators(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)

	property := vars["property"]
	table := "climres_skyspark:" + property + "_data"
	aggregation_level := vars["aggregation_level"]
	aggregation_time := vars["aggregation_time"]
	start := vars["start"]
	end := vars["end"]

	startTS, err := strconv.ParseInt(start, 10, 64)
	if err != nil {
		http.Error(w, "Invalid start timestamp", http.StatusBadRequest)
		return
	}

	endTS, err := strconv.ParseInt(end, 10, 64)
	if err != nil {
		http.Error(w, "Invalid end timestamp", http.StatusBadRequest)
		return
	}

	rows, err := hbase.GetIndoorSensorData(table, startTS, endTS)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if aggregation_level != "dwelling" {
		aggregation_level = "building"
	}

	stats := utils.ComputeIndoorIndicatorsGrouped(
		rows,
		aggregation_time,
		database.SensorBuildingMap,
		aggregation_level,
	)
	json.NewEncoder(w).Encode(stats)
}
