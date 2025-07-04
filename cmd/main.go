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
// @description This is the Climate Ready Barcelona API.
// @host localhost:8090
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

	// Conectar a MongoDB
	database.ConnectMongoDB()

	r := routes.SetupRouter()

	log.Println("Server running on port 8090")
	log.Fatal(http.ListenAndServe(":8090", r))
}
