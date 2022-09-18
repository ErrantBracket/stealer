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

	"github.com/gofiber/fiber/v2"

	"github.com/ErrantBracket/stealer/db"
	topicsModel "github.com/ErrantBracket/stealer/models"
)

/*
*
*
*/
func AddNewTopic(c *fiber.Ctx) error {
	topicsCollection := db.DB.Collection("topics")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	// Create a new Topic instance
	// The value for Topic.Topic is parsed using BodyParser
	t := new(topicsModel.Topic)
	if err := c.BodyParser(t); err != nil {
		return err
	}
	t.ID = primitive.NewObjectID()
	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()
	t.Deleted = false

	_, err := topicsCollection.InsertOne(ctx, t)
	if err != nil {
		fmt.Println(err)
	}
	
	return c.SendString(t.ID.String())
}


/*
*
*
*/
func GetTopicById(id primitive.ObjectID) *topicsModel.Topic {
	topicsCollection := db.DB.Collection("topics")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	var result *topicsModel.Topic
	err := topicsCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&result)
	if err != nil {
		return nil
	}
	fmt.Println("Topic " + result.Topic + " found (" + id.String() + ")")
	return result
}


/*
* Append the given Note ID to the given Topic references field 
* Assumes a valid Topic ID is passed in
*
*/
func AddNoteIdToTopic(topicId primitive.ObjectID, noteId primitive.ObjectID) {
	topicsCollection := db.DB.Collection("topics")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	_, err := topicsCollection.UpdateOne(ctx, 
		bson.M{"_id": topicId}, 
		bson.D{{"$push", bson.M{"references": noteId}}},
	)
	if err != nil {
		fmt.Println("Error on Update of Topic")
		fmt.Println(err)
	}
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