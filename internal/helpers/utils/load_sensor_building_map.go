package utils

import (
	"encoding/json"
	"go-api/internal/models"
	"os"
)

func LoadSensorBuildingMap(path string) (map[string]models.SensorInfo, error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var sensorMap map[string]models.SensorInfo

	decoder := json.NewDecoder(file)

	if err := decoder.Decode(&sensorMap); err != nil {
		return nil, err
	}

	return sensorMap, nil
}
