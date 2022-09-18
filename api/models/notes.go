/*
*	
*
*/
package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Note struct {
	ID			primitive.ObjectID	`bson:"_id"`
	CreatedAt	time.Time			`bson:"created_at"`
	UpdatedAt	time.Time			`bson:"updated_at"`
	Note		string				`bson:"note"`
	TopicId		primitive.ObjectID 	`bson:"topic_id"`
	Deleted		bool				`bson:"deleted"`
}

