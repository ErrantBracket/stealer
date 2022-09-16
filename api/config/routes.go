package config

import (
	notesController "github.com/ErrantBracket/stealer/controllers/notes"
	topicsController "github.com/ErrantBracket/stealer/controllers/topics"
	keywordsController "github.com/ErrantBracket/stealer/controllers/keywords"
	"github.com/gofiber/fiber/v2"
)

var api fiber.Router

func RegisterRoutes(router fiber.Router) {
	api = router.Group("api")
	registerNoteRoutes(router)
	registerTopicRoutes(router)
	registerKeywordRoutes(router)
}

// Register the routes for Notes
func registerNoteRoutes(router fiber.Router) {
	note := api.Group("/notes")
	
	note.Get("/", notesController.GetAllNotes)
	note.Post("/", notesController.CreateNote)
}

// Register the routes for Keywords
func registerKeywordRoutes(router fiber.Router) {
	keyword := api.Group("/keywords")

	keyword.Get("/", keywordsController.GetAllKeywords)
	//keyword.Patch("/", keywordsController.AddNewKeywords)
}

// Register the routes for Topics
func registerTopicRoutes(router fiber.Router) {
	topic := api.Group("/topics")

	topic.Get("/", topicsController.GetAllTopics)
}