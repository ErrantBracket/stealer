package notesModel

import (
	"time"

	"github.com/ErrantBracket/stealer/db"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var Collection *mongo.Collection

type Note struct {
	ID			primitive.ObjectID	`bson:"_id"`
	CreatedAt	time.Time			`bson:"created_at"`
	LastUpdate	time.Time			`bson:"last_update"`
	Title		string				`bson:"title"`
}

func init() {
	Collection = db.Client.Database("stealer").Collection("notes")
}