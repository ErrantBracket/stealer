package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var Ctx context.Context
var Client mongo.Client

func ConnectDb(url string) {
	Client, err := mongo.NewClient(options.Client().ApplyURI(url))
	if err != nil {
		log.Fatal(err)
	}

	Ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	err = Client.Connect(Ctx)
	if err != nil {
		log.Fatal(err)
	}

	defer cancel()

	err = Client.Ping(Ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
}
