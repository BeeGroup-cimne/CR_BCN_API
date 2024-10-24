package database

import (
	"context"
	"go-api/internal/helpers/env"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client

// ConnectMongoDB
func ConnectMongoDB() {
	uri := env.GetVar("MONGODB_URI")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("Error conectando a MongoDB: ", err)
	}

	// Verificar la conexión
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("Error al hacer ping a MongoDB: ", err)
	}

	MongoClient = client
	log.Println("Conectado a MongoDB!")
}
