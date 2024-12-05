package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

func GetPort() string {
	// load env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// fetch DB uri frm env file
	port := os.Getenv("DB_PORT")
	if port == "" {
		log.Fatal("DB_PORT not set in .env file")
	}

	return port
}
