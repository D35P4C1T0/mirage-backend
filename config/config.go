package config

import (
	"github.com/joho/godotenv"
	"log"
)

func init() {
	// load env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
