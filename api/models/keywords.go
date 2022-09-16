package models

import (
	"time"

	
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var KeywordsCollection *mongo.Collection

type Keyword struct {
	ID			primitive.ObjectID		`bson:"_id"`
	Keyword		string					`bson:"keyword"`
	CreatedAt	time.Time				`bson:"create_at"`	
	UpdatedAt	time.Time				`bson:"updated_at"`
	References	[]primitive.ObjectID	`bson:"references"`
}
