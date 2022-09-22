/*
*
*
 */
package topicsController

import (
	"context"
	"errors"
	"fmt"
	"log"
	"reflect"
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
	fmt.Println("TopicId " + id.String() + " found")
	return true
}


/*
*
*
*/
func GetNextSequence(id primitive.ObjectID) (error, int) {
	topicsCollection := db.DB.Collection("topics")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	//var result *topicsModel.TopicSequenceOnly

	matchStage   := bson.D{{"$match", bson.D{{"_id", id}}}}
	unwindStage  := bson.D{{"$unwind", "$notes"}}
	replaceRoot  := bson.D{{"$replaceRoot", bson.D{{"newRoot", "$notes"}}}}
	sortStage    := bson.D{{"$sort", bson.M{"sequence": -1}}}
	limitStage   := bson.D{{"$limit", 1}}
	projectStage := bson.D{{"$project", bson.M{"_id": 0, "sequence": 1}}}	
	
	cur, err := topicsCollection.Aggregate(ctx, mongo.Pipeline{matchStage, unwindStage, replaceRoot, sortStage, limitStage, projectStage})
	if err != nil {
		log.Fatal(err)
	}
	

	var sequenceInfo []bson.M

	if err = cur.All(ctx, &sequenceInfo); err != nil {
		log.Fatal(err)
	}

	fmt.Println(sequenceInfo)
	iter := reflect.ValueOf(sequenceInfo).MapRange()
	for iter.Next() {
		k := iter.Key().Interface()
		v := iter.Value().Interface()
		fmt.Println(k, v)
	}

	err = errors.New("Error")
	return err, 0
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