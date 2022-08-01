package notesController

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ErrantBracket/stealer/db"
	notesModel "github.com/ErrantBracket/stealer/models"
	keywordsController "github.com/ErrantBracket/stealer/controllers/keywords"
	"github.com/gofiber/fiber/v2"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const keywordStart string = "{{"
const keywordEnd   string = "}}"

type Note struct {
	Note string `json:"note"`
}

func GetAllNotes(c *fiber.Ctx) error {
	filter := bson.D{{}}
	notes, err := filterNotes(filter)
	if err != nil {

	}
	return c.SendString("PUT: Hello World " + strconv.Itoa(len(notes)))
}

func CreateNote(c *fiber.Ctx) error {
	//_, err := notes.Collection.InsertOne(db.Ctx, note)
	n := new(Note)
	if err := c.BodyParser(n); err != nil {
		return err
	}
	sanitiseNote(n)
	return c.SendString("PUT: Hello World")
}


func filterNotes(filter interface{}) ([]*notesModel.Note, error) {
	var notes []*notesModel.Note

	cur, err := notesModel.NotesCollection.Find(db.Ctx, filter)
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

func sanitiseNote(note *Note) string {

	// Check if there is even a possible keyword contained in <note>
	if strings.Index(note.Note, keywordStart) == -1 || strings.Index(note.Note, keywordEnd) == -1 {
		return note.Note
	}

	var keywords []string
	var cleanString string

	// We must have at least one keyword template match
	//kws := strings.Index(note, keywordStart)
	for kws := strings.Index(note.Note, keywordStart); kws != -1; kws =  strings.Index(note.Note, keywordStart) {
		cleanString = cleanString + note.Note[:kws]
		note.Note = note.Note[kws+2:]
		kwe := strings.Index(note.Note, keywordEnd)
		keywords = append(keywords, note.Note[:kwe])
		note.Note = note.Note[kwe+2:]	
	}
	cleanString = cleanString + note.Note
	fmt.Println(keywords)
	fmt.Println(cleanString)
	keywordsController.AddNewKeywords(keywords)
	return cleanString
}