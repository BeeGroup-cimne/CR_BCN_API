package main

import (
	"go-api/internal/database"
	"go-api/internal/routes"
	"log"
	"net/http"
)

func main() {

	// Conectar a MongoDB
	database.ConnectMongoDB()

	r := routes.SetupRouter()

	log.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
