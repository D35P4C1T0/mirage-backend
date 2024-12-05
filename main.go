package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"mirage-backend/config"
	"mirage-backend/database"
	"mirage-backend/routes"
)

func main() {
	database.SetupDatabase()

	router := gin.Default()
	router.Use(cors.Default())
	routes.InitRoutes(router)

	port := config.GetPort()
	if port == "" {
		port = "8080"
	}

	err := router.Run(":" + port)
	if err != nil {
		log.Fatalln(err)
	}
}
