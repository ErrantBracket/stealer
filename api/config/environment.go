package config

import (
	"os"
	"errors"
)

var Port string
var DbUrl string

func GetEnvironment() (error) {
	var exists bool

	Port, exists = os.LookupEnv("PORT"); if !exists {
		return errors.New("PORT environment variable not set")
	}

	DbUrl, exists = os.LookupEnv("DB_URL"); if !exists {
		return errors.New("DB_URL environment variable not set")
	}
	return nil
}

// mongodb://mongo:27017/note-taker
