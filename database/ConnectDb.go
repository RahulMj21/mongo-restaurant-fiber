package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DbName = "mongo-restaurant-fiber"
var DbUrl = "mongodb://localhost:27017/" + DbName

func DbInstance() *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI(DbUrl))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("DB Connected..")

	return client
}

var Client = DbInstance()

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	return client.Database(DbName).Collection(collectionName)
}
