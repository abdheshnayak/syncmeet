package db

import (
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func NewDB(url string) *mongo.Client {
	client, err := mongo.Connect(options.Client().ApplyURI(url))
	if err != nil {
		panic(err)
	}

	return client
}
