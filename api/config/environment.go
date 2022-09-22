package config

import (
	"os"
)

var Port string
var DbUrl string

func GetEnvironment() (error) {
	var exists bool

	Port, exists = os.LookupEnv("PORT"); if !exists {
		Port = "3000"
	}

	DbUrl, exists = os.LookupEnv("DB_URL"); if !exists {
		DbUrl = "mongodb://localhost:27017"
	}
	return nil
}
