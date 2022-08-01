package models

import (
	"github.com/ErrantBracket/stealer/db"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var KeywordsCollection *mongo.Collection

type Keyword struct {
	ID			primitive.ObjectID		`bson:"_id"`
	Keyword		string					`bson:"keyword"`
	References	[]primitive.ObjectID	`bson:"_id"`
}

func init() {
	KeywordsCollection = db.Client.Database("stealer").Collection("keywords")
}