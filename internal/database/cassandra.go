package database

import (
	"go-api/internal/helpers/utils"
	"log"
	"strings"
	"time"

	"github.com/gocql/gocql"
	"go-api/internal/helpers/env"
)

var CassandraSession *gocql.Session

// BuildingAggregationDwellings is a fixed (config-driven, not caller-supplied)
// list of dwelling identifiers to include when aggregating at "building"
// level. Set via the BUILDING_AGGREGATION_DWELLINGS env var as a
// comma-separated list (matching the raw dwelling token as it appears in
// the Cassandra hash, e.g. "3-1,4-2,3-2h"). Empty/unset means "no filter,
// include every dwelling".
var BuildingAggregationDwellings []string

// ConnectCassandraDB opens a session against the Cassandra cluster and
// keeps a single shared *gocql.Session (gocql sessions are safe for
// concurrent use, so no pooling logic is needed on top of it).
func ConnectCassandraDB() {
	hosts := strings.Split(env.GetVar("CASSANDRA_HOSTS"), ",")
	keyspace := env.GetVar("CASSANDRA_KEYSPACE")
	user := env.GetVar("CASSANDRA_USER")
	password := env.GetVar("CASSANDRA_PASSWORD")

	cluster := gocql.NewCluster(hosts...)
	cluster.Keyspace = keyspace
	cluster.Consistency = gocql.LocalQuorum
	cluster.Timeout = 15 * time.Second
	cluster.ConnectTimeout = 10 * time.Second

	if user != "" {
		cluster.Authenticator = gocql.PasswordAuthenticator{
			Username: user,
			Password: password,
		}
	}

	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatal("Error connecting to Cassandra: ", err)
	}

	CassandraSession = session
	log.Println("Connect to Cassandra!")

	SensorBuildingMap, err = utils.LoadSensorBuildingMap("internal/data/sensor_building_map_cassandra.json")
	if err != nil {
		log.Fatal("Error loading sensor building map: ", err)
	}
	log.Printf("Loaded %d sensors into SensorBuildingMap.\n", len(SensorBuildingMap), SensorBuildingMap)

	//SensorBuildingMap, err = utils.LoadSensorBuildingMapFromCassandra(session, "skyspark.climres_static_pt15m")
	//if err != nil {
	//	log.Fatal("Error loading sensor building map from Cassandra: ", err)
	//}
	//log.Printf("Loaded %d sensors into SensorBuildingMap%v.\n", len(SensorBuildingMap), SensorBuildingMap)

	//if raw := env.GetVar("BUILDING_AGGREGATION_DWELLINGS"); raw != "" {
	//	for _, d := range strings.Split(raw, ",") {
	//		d = strings.TrimSpace(d)
	//		if d != "" {
	//			BuildingAggregationDwellings = append(BuildingAggregationDwellings, d)
	//		}
	//	}
	//}
	//log.Printf("Building-level aggregation restricted to %d dwelling(s): %v\n",
	//	len(BuildingAggregationDwellings), BuildingAggregationDwellings)
}
