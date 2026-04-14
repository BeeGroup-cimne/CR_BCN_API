package hbase

import (
	"context"
	"fmt"
	"go-api/gen-go/hbase"
	"go-api/internal/database"
	"go-api/internal/models"
	"strconv"
)

func GetIndoorSensorData(table string, startTS int64, endTS int64) ([]models.SensorRow, error) {

	err := database.HbaseTransport.Open()
	if err != nil {
		return nil, err
	}
	defer database.HbaseTransport.Close()

	startRow := fmt.Sprintf("%d", startTS)
	stopRow := fmt.Sprintf("%d~", endTS)

	scanner, err := database.HbaseClient.ScannerOpenWithStop(
		context.Background(),
		hbase.Text(table),
		hbase.Text(startRow),
		hbase.Text(stopRow),
		nil,
		nil,
	)
	if err != nil {
		return nil, err
	}

	defer database.HbaseClient.ScannerClose(context.Background(), scanner)

	result := make([]models.SensorRow, 0, 10000)

	for {

		rows, err := database.HbaseClient.ScannerGetList(
			context.Background(),
			scanner,
			10000,
		)
		if err != nil {
			return nil, err
		}

		if len(rows) == 0 {
			break
		}

		for _, row := range rows {

			rowKey := string(row.Row)

			for _, col := range row.Columns {

				v, err := strconv.ParseFloat(string(col.Value), 64)
				if err != nil {
					continue
				}

				result = append(result, models.SensorRow{
					RowKey: rowKey,
					Value:  v,
				})
			}
		}
	}

	return result, nil
}
