package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var Client *mongo.Client
var DB *mongo.Database

func ConnectDb(url string) {
	var err error
	
	Client, err = mongo.NewClient(options.Client().ApplyURI(url))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	err = Client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	defer cancel()

	err = Client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	DB = Client.Database("stealer")

	fmt.Println("Connected to DB")
}


