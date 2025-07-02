package main

import (
	"github.com/joho/godotenv"
	_ "go-api/docs"
	"go-api/internal/database"
	"go-api/internal/routes"
	"log"
	"net/http"
)

// @title (CRAPI) Climate Ready Barcelona API
// @version 1.0
// @description This is an example API with Swagger documentation using mux.
// @host localhost:8090
// @BasePath /
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	// Conectar a MongoDB
	database.ConnectMongoDB()

	r := routes.SetupRouter()

	log.Println("Server running on port 8090")
	log.Fatal(http.ListenAndServe(":8090", r))
}
