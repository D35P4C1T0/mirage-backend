package routes

import (
	"github.com/gin-gonic/gin"
	"mirage-backend/controllers"
)

func SetupUserRoutes(api *gin.RouterGroup) {
	controllers.InitializeUserController()

	userRoutes := api.Group("/users")
	{
		// Get User Profile
		userRoutes.GET("/:userId", controllers.GetUserProfile)

		// Update User Profile
		userRoutes.PUT("/:userId", controllers.UpdateUserProfile)

		// Delete User
		userRoutes.DELETE("/:userId", controllers.DeleteUser)

		// Get All Users
		userRoutes.GET("/", controllers.GetAllUsers)
	}
}
