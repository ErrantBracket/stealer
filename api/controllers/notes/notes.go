package notesController

import (
	"strconv"

	"github.com/ErrantBracket/stealer/db"
	notesModel "github.com/ErrantBracket/stealer/models"
	"github.com/gofiber/fiber/v2"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetAllNotes(c *fiber.Ctx) error {
	filter := bson.D{{}}
	notes, err := filterNotes(filter)
	if err != nil {

	}
	return c.SendString("PUT: Hello World " + strconv.Itoa(len(notes)))
}

func CreateNote(c *fiber.Ctx) error {
	//_, err := notes.Collection.InsertOne(db.Ctx, note)
	return c.SendString("PUT: Hello World")
}


func filterNotes(filter interface{}) ([]*notesModel.Note, error) {
	var notes []*notesModel.Note

	cur, err := notesModel.Collection.Find(db.Ctx, filter)
	if err != nil {
		return notes, err
	}

	// Iterate over cursor returned by Collection.Find
	// and decode each document into an instance of Note.
	// Each Note is then appended to the slice of Notes
	for cur.Next(db.Ctx) {
		var n notesModel.Note
		err := cur.Decode(&n)
		if err != nil {
			return notes, err
		}
		notes = append(notes, &n)
	}
	cur.Close(db.Ctx)

	if len(notes) == 0 {
		return notes, mongo.ErrNoDocuments
	}

	return notes, nil
}