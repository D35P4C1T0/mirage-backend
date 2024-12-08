package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"mirage-backend/config"
	"mirage-backend/database"
	"mirage-backend/routes"
	"time"
)

func main() {
	database.SetupDatabase()

	router := gin.Default()
	//router.Use(cors.Default())

	router.Use(cors.New(cors.Config{
		//AllowOrigins: []string{"http://172.0.0.1:5500"},                            // Allow all origins
		AllowMethods:    []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}, // Allow all methods
		AllowAllOrigins: true,
		AllowHeaders:    []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:   []string{"Content-Length"}, AllowCredentials: true, MaxAge: 12 * time.Hour}))

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
