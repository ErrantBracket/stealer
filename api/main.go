package main

import (
	"github.com/ErrantBracket/stealer/db"
	"github.com/ErrantBracket/stealer/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)


func main() {
	config.GetEnvironment()
	db.ConnectDb(config.DbUrl)

	app := fiber.New()
	app.Use(logger.New())
	// Run the route definitions in config > routes.go
	config.RegisterRoutes(app)

	app.Listen(":" + config.Port)
}
