package routes

import (
	"github.com/gin-gonic/gin"
	"mirage-backend/controllers"
)

func SetupProfilePictureRoutes(api *gin.RouterGroup) {
	// Route group for profile picture-related operations
	profilePictureRoutes := api.Group("/profilepictures")
	{
		// Upload profile picture
		profilePictureRoutes.POST("/user/:userId", controllers.UploadProfilePicture)

		//// Get all profile pictures
		//profilePictureRoutes.GET("/", controllers.GetAllProfilePictures)
		//
		//// Get profile picture by user ID
		//profilePictureRoutes.GET("/user/:userId", controllers.GetProfilePictureByUserID)
		//
		//// Get profile picture by ID
		//profilePictureRoutes.GET("/:profilePictureId", controllers.GetProfilePictureByID)
		//
		//// Update profile picture by ID
		//profilePictureRoutes.PUT("/:profilePictureId", controllers.UpdateProfilePicture)
		//
		//// Delete profile picture by ID
		//profilePictureRoutes.DELETE("/:profilePictureId", controllers.DeleteProfilePicture)
	}
}
