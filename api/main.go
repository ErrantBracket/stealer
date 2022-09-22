package main

import (
	"log"

	"github.com/ErrantBracket/stealer/config"
	"github.com/ErrantBracket/stealer/db"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)


func main() {
	err := config.GetEnvironment()
	if err != nil {
		
	}
	
	app := fiber.New()
	app.Use(logger.New())

	db.ConnectDb(config.DbUrl)
	
	// Run the route definitions in config > routes.go
	config.RegisterRoutes(app)

	err = app.Listen(":" + config.Port)
	if err != nil {
		log.Fatal("App failed to start")
	}
}
