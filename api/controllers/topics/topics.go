/*
*
*
 */
package topicsController

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/gofiber/fiber/v2"

	"github.com/ErrantBracket/stealer/db"
	notesModel "github.com/ErrantBracket/stealer/models"
	topicsModel "github.com/ErrantBracket/stealer/models"
)

/*
*
*
 */
func CreateNewTopic(c *fiber.Ctx) error {
	topicsCollection := db.DB.Collection("topics")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	// Create a new Topic instance
	// The value for Topic.Topic is parsed using BodyParser
	t := new(topicsModel.Topic)
	if err := c.BodyParser(t); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	t.ID = primitive.NewObjectID()
	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()
	t.Deleted = false
	t.Keywords = make([]string, 0)
	t.Notes = make([]notesModel.Note, 0)

	_, err := topicsCollection.InsertOne(ctx, t)
	if err != nil {
		fmt.Println(err)
	}
	
	// Return the new topic _id
	return c.SendString(t.ID.String())
}


/*
*
*
*/
func GetTopicById(id primitive.ObjectID) (error, *topicsModel.Topic) {
	fmt.Println("GetTopicById(" + id.String() + ")")
	topicsCollection := db.DB.Collection("topics")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	var result *topicsModel.Topic
	err := topicsCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&result)
	if err != nil {
		return err, nil
	}
	fmt.Println("Topic " + result.Topic + " found (" + id.String() + ")")
	return nil, result
}


/*
* Append the given Note to the given Topic references field 
* Assumes a valid Topic ID is passed in
*
*/
func AddNoteToTopic(topicId primitive.ObjectID, note *notesModel.Note) error {
	topicsCollection := db.DB.Collection("topics")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	_, err := topicsCollection.UpdateOne(ctx, 
		bson.M{"_id": topicId}, 
		bson.D{{"$push", bson.M{"notes": note}}},
	)
	if err != nil {
		return err
	}
	return nil
}


/*
* Check the given Topic ID is valid 
*
* @param	id		Mongo Object ID to check in Topics collection
* @return   bool	
*/
func IsTopicValid(id primitive.ObjectID) bool {
	topicsCollection := db.DB.Collection("topics")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	var result *topicsModel.TopicIdOnly
	projection := bson.D{{"_id", 1}}
	opts := options.FindOne().SetProjection(projection)
	err := topicsCollection.FindOne(ctx, bson.M{"_id": id}, opts).Decode(&result)
	if err != nil {
		return false
	}
	//fmt.Println("TopicId " + id.String() + " found")
	return true
}


/*
* Get the highest sequence number + 1
* Honestly this took me far too long to figure out and seems quite messy
*
* Based on the following Mongosh:
* db.topics.aggregate( [ { $match:{"_id": ObjectId("6328d8083c35cf86c59bcdf1")}}, 
*		{$unwind:"$notes"}, 
*		{$set: {sequence: "$notes.sequence"}}, 
*		{$sort:{"sequence":-1}}, 
*		{$limit: 1}, 
*		{$project: {"_id":0, "sequence":1}} ])
*
*
* Get a count of the number of notes subdocuments
* - If zero, then return 1 as the next sequence
* 
*/
func GetNextSequence(id primitive.ObjectID) (int, error) {
	num, err := GetNoteCount(id)
	if err != nil {
		return 0, err
	}

	if num == 0 {
		return 1, nil
	}

	topicsCollection := db.DB.Collection("topics")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	// Get the topic with a matching _id
	matchStage   := bson.D{{"$match", bson.D{{"_id", id}}}}
	// Deconstruct an array field to a document for each instance of Notes
	unwindStage  := bson.D{{"$unwind", "$notes"}}
	// Create a new field copied from the subdocument
	setStage 	 := bson.D{{"$set", bson.M{"sequence": "$notes.sequence"}}}
	// Sort by sequence (descending)
	sortStage    := bson.D{{"$sort", bson.M{"sequence": -1}}}
	// Give me 1 document only
	limitStage   := bson.D{{"$limit", 1}}
	// Give me the sequence field (_id is returned by default so we supress it)
	projectStage := bson.D{{"$project", bson.M{"_id": 0, "sequence": 1}}}	
	
	cur, err := topicsCollection.Aggregate(ctx, mongo.Pipeline{matchStage, unwindStage, setStage, sortStage, limitStage, projectStage})
	if err != nil {
		log.Fatal(err)
	}

	type Sequence2 struct {
		Sequence	int			`bson:"sequence"`
	}
	type SequenceInfo struct {
		Sequence 	Sequence2 	`bson:"inline,omitempty"`
	}

	//var sequenceInfo []bson.M
	// [{"sequence": {"$numberInt":"4"}}]
	// [map[sequence:4]]

	var sequenceInfo []SequenceInfo
	if err = cur.All(ctx, &sequenceInfo); err != nil {
		log.Fatal(err)
	}
	//fmt.Println(sequenceInfo)  // [{{4}}]

	seq := sequenceInfo[0].Sequence.Sequence
	return seq + 1, err
}

/*
*
* Mongosh command:
* db.topics.aggregate( [ {$match:{"_id": ObjectId("6328d8083c35cf86c59bcdf1")}}, 
*		{$project: {numberOfNotes: {$size: "$notes"}}} ] )
* 
*/
func GetNoteCount(id primitive.ObjectID) (int, error) {
	topicsCollection := db.DB.Collection("topics")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	// Get the topic with a matching _id
	matchStage   := bson.D{{"$match", bson.D{{"_id", id}}}}
	// Give me the sequence field (_id is returned by default so we supress it)
	projectStage := bson.D{{"$project", bson.D{{"numberOfNotes", bson.M{"$size": "$notes"}}}}}
	
	cur, err := topicsCollection.Aggregate(ctx, mongo.Pipeline{matchStage, projectStage})
	if err != nil {
		log.Fatal(err)
	}

	type NoteCount struct {
		Id 		primitive.ObjectID	`bson:"_id"`
		Num		int					`bson:"numberOfNotes"`
	}

	var noteCount []NoteCount
	if err = cur.All(ctx, &noteCount); err != nil {
		log.Fatal(err)
	}

	num := noteCount[0].Num
	
	return num, nil
}

/*
*
*
*/
func GetAllTopics(c *fiber.Ctx) error {
	topicsCollection := db.DB.Collection("topics")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	var results []*topicsModel.Topic
	cursor, err := topicsCollection.Find(ctx, bson.M{"deleted": bson.M{"$eq": false}})
	defer cursor.Close(ctx)
	if err != nil {
		panic(err)
	}
	for cursor.Next(ctx) {
		var result = new(topicsModel.Topic)
		if err := cursor.Decode(&result); err != nil {
			log.Fatal(err)
		}
		results = append(results, result)
	}
	fmt.Println(results)

	c.JSON(results)
	return c.SendStatus(200)
}