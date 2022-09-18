/*
*	
*
*/
package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Topic struct {
	ID			primitive.ObjectID		`bson:"_id"`
	CreatedAt	time.Time				`bson:"created_at"`
	UpdatedAt	time.Time				`bson:"updated_at"`
	Topic		string					`bson:"topic"`
	References 	[]primitive.ObjectID	`bson:"references"`
	Deleted		bool					`bson:"deleted"`
}
