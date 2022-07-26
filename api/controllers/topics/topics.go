package topicsController

import "github.com/gofiber/fiber/v2"


func GetAllTopics(c *fiber.Ctx) error {
	return c.SendString("Hello World")
}