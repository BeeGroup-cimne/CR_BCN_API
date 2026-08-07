package main

import (
	"github.com/joho/godotenv"
	_ "go-api/docs"
	"go-api/internal/database"
	"go-api/internal/routes"
	"net/http"

	//"go-api/internal/database"

	"log"
)

// @title 🚀 CRAPI – Climate Ready Barcelona API
// @version 1.0
// @description 🌍 Climate Ready BCN – Smart Data for a Resilient Future
// @host crbcn-go-api.kubernetes.tech.beegroup-cimne.com
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Provide your access token like this: "Bearer {token}"
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	database.ConnectMongoDB()
	database.ConnectCassandraDB()
	//database.ConnectHBaseDB()

	r := routes.SetupRouter()

	//log.Println("Server running on port 8090")
	log.Fatal(http.ListenAndServe(":8090", r))
}
