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
	Notes	 	[]Note					`bson:"notes"`
	Keywords	[]string				`bson:"keywords"`
	Deleted		bool					`bson:"deleted"`
}

type TopicIdOnly struct {
	ID 			primitive.ObjectID		`bson:"_id"`
}

type TopicSequenceOnly struct {
	Sequence 	int						`bson:"sequence"`
}
