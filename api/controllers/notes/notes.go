/*
*	
*
*/
package notesController

import (
	"context"
	"fmt"
	"strings"
	"strconv"
	"time"

	"github.com/ErrantBracket/stealer/db"
	keywordsModel 	"github.com/ErrantBracket/stealer/models"
	notesModel 		"github.com/ErrantBracket/stealer/models"
	keywordsController 	"github.com/ErrantBracket/stealer/controllers/keywords"
	topicsController 	"github.com/ErrantBracket/stealer/controllers/topics"
	"github.com/gofiber/fiber/v2"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetAllNotes(c *fiber.Ctx) error {
	filter := bson.D{{}}
	notes, err := filterNotes(filter)
	if err != nil {

	}
	return c.SendString("PUT: Hello World " + strconv.Itoa(len(notes)))
}

/*
* Create a new Note instance
* POST	/notes
* 		TopicId:	String representation of Mongo _id for related topic
*		Note:		String
*
* Check that topic exists
* Insert note
* Parse note for keywords
* return _id for new note entry 
*/
func CreateNote(c *fiber.Ctx) error {
	notesCollection := db.DB.Collection("notes")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	// Create a new Note instance
	// The value for Note.TopidId, and Note.Note is parsed using BodyParser
	n := new(notesModel.Note)
	if err := c.BodyParser(n); err != nil {
		return err
	}
	n.ID = primitive.NewObjectID()
	n.CreatedAt = time.Now()
	n.UpdatedAt = time.Now()
	n.Deleted = false

	// Check Topic is valid and append note ID to reference
	topic := topicsController.GetTopicById(n.TopicId)
	if topic == nil {
		return c.SendStatus(404)
	} else {
		topicsController.AddNoteIdToTopic(n.TopicId, n.ID)
	}

	result, err := notesCollection.InsertOne(ctx, n)
	if err != nil {
		fmt.Println(err)
	} else {
		checkNoteForKeywords(n.Note, n.ID)
		fmt.Println(result)
	}
	
	return c.SendString(n.ID.String())
}


func filterNotes(filter interface{}) ([]*notesModel.Note, error) {
	var notes []*notesModel.Note

	notesCollection := db.DB.Collection("notes")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	cur, err := notesCollection.Find(ctx, filter)
	if err != nil {
		return notes, err
	}

	// Iterate over cursor returned by Collection.Find
	// and decode each document into an instance of Note.
	// Each Note is then appended to the slice of Notes
	for cur.Next(ctx) {
		var n notesModel.Note
		err := cur.Decode(&n)
		if err != nil {
			return notes, err
		}
		notes = append(notes, &n)
	}
	cur.Close(ctx)

	if len(notes) == 0 {
		return notes, mongo.ErrNoDocuments
	}

	return notes, nil
}

/*
* Parse a string for keywords
* Keywords are added to a slice and passed to be added to the database
*/
func checkNoteForKeywords(note string, id primitive.ObjectID) {
	if strings.Contains(note, keywordsModel.KeywordStart) && strings.Contains(note, keywordsModel.KeywordEnd) {
		var keywords []string
		
		// Loop over the string while we find the start of a keyword marker
		for kws := strings.Index(note, keywordsModel.KeywordStart); kws != -1; kws =  strings.Index(note, keywordsModel.KeywordStart) {
			note = note[kws+2:]
			kwe := strings.Index(note, keywordsModel.KeywordEnd)
			keywords = append(keywords, note[:kwe])
			note = note[2:]	
		}
		fmt.Println(keywords)
		
		// If we have any keyword entries then add them
		if len(keywords) > 0 {
			keywordsController.AddNewKeywords(keywords, id)
		}
	}
}