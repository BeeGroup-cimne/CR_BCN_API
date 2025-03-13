package main

import (
	"go-api/internal/database"
	"go-api/internal/routes"
	"log"
	"net/http"
)

// @title My API
// @version 1.0
// @description This is an example API with Swagger documentation using mux.
// @host localhost:8080
// @BasePath /api/v1
func main() {

	// Conectar a MongoDB
	database.ConnectMongoDB()

	r := routes.SetupRouter()

	log.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
