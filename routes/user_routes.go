package routes

import (
	"github.com/gin-gonic/gin"
	"mirage-backend/controllers"
)

// SetupUserRoutes sets up all routes related to User
func SetupUserRoutes(rg *gin.RouterGroup) {
	rg.POST("/users", controllers.CreateUser)
	// Add other routes here
}
