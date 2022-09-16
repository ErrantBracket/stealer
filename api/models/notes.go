package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var NotesCollection *mongo.Collection

type Note struct {
	ID			primitive.ObjectID	`bson:"_id"`
	CreatedAt	time.Time			`bson:"created_at"`
	LastUpdate	time.Time			`bson:"last_update"`
	Title		string				`bson:"title"`
}

