package database

import (
	"fmt"
	"github.com/apache/thrift/lib/go/thrift"
	"go-api/gen-go/hbase"
	"go-api/internal/helpers/env"
	"go-api/internal/helpers/utils"
	"go-api/internal/models"
	"log"
)

var HbaseClient *hbase.HbaseClient
var HbaseTransport *thrift.TBufferedTransport
var SensorBuildingMap map[string]models.SensorInfo

func ConnectHBaseDB() {
	hbaseHost := env.GetVar("HBASE_HOST")
	hbasePort := env.GetVar("HBASE_PORT")
	addr := fmt.Sprintf("%s:%s", hbaseHost, hbasePort)
	socket := thrift.NewTSocketConf(addr, &thrift.TConfiguration{})
	transport := thrift.NewTBufferedTransport(socket, 8192)
	protocol := thrift.NewTBinaryProtocolConf(transport, nil)
	client := hbase.NewHbaseClientProtocol(transport, protocol, protocol)
	HbaseClient = client
	HbaseTransport = transport

	log.Println("Connect to HBase!")

	var err error
	SensorBuildingMap, err = utils.LoadSensorBuildingMap("internal/data/sensor_building_map.json")

	if err != nil {
		panic(err)
	}
}
