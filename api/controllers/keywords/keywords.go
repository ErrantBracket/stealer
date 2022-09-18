/*
*
*
*/
package keywordsController

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/ErrantBracket/stealer/db"
	keywordsModel "github.com/ErrantBracket/stealer/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


func GetAllKeywords(c *fiber.Ctx) error {
	return c.SendString("Hello World")
}


func AddNewKeywords(kws []string, id primitive.ObjectID) error {	
	keywordsCollection := db.DB.Collection("keywords")
	ctx, _ := context.WithTimeout(context.Background(), 10 * time.Second)


	for _, kw := range kws {
		fmt.Println("Keyword: " + kw)
	
		// Check if the keyword already exists
		// - if so, update the existing one with UpdatedAt and added Reference
		// - if not, add a new entry
		result := new(keywordsModel.Keyword)
		err := keywordsCollection.FindOne(ctx, bson.M{"keyword": kw}).Decode(&result)
		if err == nil {
			// No error from the FineOne so we have a document to update
			_, err := keywordsCollection.UpdateOne(ctx, 
				bson.M{"_id": result.ID}, 
				bson.D{{"$set", bson.M{"updated_at": time.Now(), "references": append(result.References, id)}}},
			)
			if err != nil {
				fmt.Println("Error on Update on")
				fmt.Println(err)
			}
		} else {
			if err == mongo.ErrNoDocuments {
				k := new(keywordsModel.Keyword)
				k.ID = primitive.NewObjectID()
				k.Keyword = kw
				k.CreatedAt = time.Now()
				k.UpdatedAt = time.Now()
				k.References = []primitive.ObjectID{id}
				result, err := keywordsCollection.InsertOne(ctx, k)
				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Println(result)
				}
			}
		}
	}

	return nil
}
