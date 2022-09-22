/*
*	
*
*/
package models

import (
	"time"
)

type Note struct {
	CreatedAt	time.Time			`bson:"created_at"`
	UpdatedAt	time.Time			`bson:"updated_at"`
	Note		string				`bson:"note"`
	Deleted		bool				`bson:"deleted"`
	Sequence	int					`bson:"sequence"`
}

