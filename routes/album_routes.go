package routes

import (
	"github.com/gin-gonic/gin"
	"mirage-backend/controllers"
)

// SetupAlbumRoutes sets up all routes related to Album
func SetupAlbumRoutes(api *gin.RouterGroup) {

	//controllers.InitializeAlbumController()

	// Route group for album-related picture operations
	albumRoutes := api.Group("/albums")
	{
		// Create a new album
		albumRoutes.POST("/", controllers.CreateAlbum)

		// Get all albums
		albumRoutes.GET("/", controllers.GetAllAlbums)

		// Get a specific album by ID
		albumRoutes.GET("/:albumId", controllers.GetAlbumByID)
	}
}
