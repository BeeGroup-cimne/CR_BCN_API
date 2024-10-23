package models

import "go.mongodb.org/mongo-driver/bson"

type MongoMessage struct {
	Status  string `json:"status"`
	Message bson.M `json:"message"`
}
