package models

import (
	"time"

	
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const KeywordStart string = "{{"
const KeywordEnd   string = "}}"

type Keyword struct {
	ID			primitive.ObjectID		`bson:"_id"`
	Keyword		string					`bson:"keyword"`
	CreatedAt	time.Time				`bson:"create_at"`	
	UpdatedAt	time.Time				`bson:"updated_at"`
	References	[]primitive.ObjectID	`bson:"references"`
}
