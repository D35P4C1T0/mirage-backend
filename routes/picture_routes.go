package routes

import (
	"github.com/gin-gonic/gin"
	"mirage-backend/controllers"
)

func SetupPictureRoutes(api *gin.RouterGroup) {

	controllers.InitializePictureController()
	// Route group for album-related picture operations
	albumRoutes := api.Group("/albums/:albumId/pictures")
	{
		// Upload picture to a specific album
		albumRoutes.POST("/", controllers.UploadPictureToAlbum)

		// Get all pictures in a specific album
		albumRoutes.GET("/", controllers.GetPicturesInAlbum)
	}

	// Route group for picture-related operations
	pictureRoutes := api.Group("/pictures")
	{
		// Get all pictures
		pictureRoutes.GET("/", controllers.GetAllPictures)

		// Singular picture operations
		pictureRoutes.GET("/:pictureId", controllers.GetPictureByID)   // Retrieve a specific picture by ID
		pictureRoutes.DELETE("/:pictureId", controllers.DeletePicture) // Delete a specific picture by ID
	}
}
