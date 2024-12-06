package config

import (
	"log"
	"os"
)

func GetDatabaseName() string {
	// fetch DB uri frm env file
	databaseName := os.Getenv("DB_DATABASE")
	if databaseName == "" {
		log.Fatal("DB_DATABASE not set in .env file")
	}

	return databaseName
}

func GetDatabaseURI() string {
	// fetch DB uri frm env file
	databaseURI := os.Getenv("DB_URI")
	if databaseURI == "" {
		log.Fatal("DB_URI not set in .env file")
	}

	return databaseURI
}

func GetPort() string {
	// fetch DB uri frm env file
	port := os.Getenv("BACKEND_PORT")
	if port == "" {
		log.Fatal("BACKEND_PORT not set in .env file")
	}

	return port
}
