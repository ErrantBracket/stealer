package keywordsController

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/ErrantBracket/stealer/db"
	keywordsModel "github.com/ErrantBracket/stealer/models"

	//"go.mongodb.org/mongo-driver/bson"
	//"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const KeywordStart string = "{{"
const KeywordEnd   string = "}}"

func GetAllKeywords(c *fiber.Ctx) error {
	return c.SendString("Hello World")
}


func AddNewKeywords(kws []string) error {	
	keywordsCollection := db.DB.Collection("keywords")
	ctx, _ := context.WithTimeout(context.Background(), 10 * time.Second)

	for _, kw := range kws {
		fmt.Println("Keyword: " + kw)
		k := new(keywordsModel.Keyword)
		k.ID = primitive.NewObjectID()
		k.Keyword = kw
		k.CreatedAt = time.Now()
		k.UpdatedAt = time.Now()
		
		
		result, err := keywordsCollection.InsertOne(ctx, k)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(result)
		}
	}

	return nil
}
