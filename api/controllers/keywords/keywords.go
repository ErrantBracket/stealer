package keywordsController

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
)

func GetAllKeywords(c *fiber.Ctx) error {
	return c.SendString("Hello World")
}


func AddNewKeywords(kws []string) error {
	for _, kw := range kws {
		fmt.Println("Keyword: " + kw)
	}


	return nil
}
