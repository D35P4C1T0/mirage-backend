package routes

import (
	"github.com/gin-gonic/gin"
	"mirage-backend/controllers"
)

// SetupAlbumRoutes sets up all routes related to Album
func SetupAlbumRoutes(rg *gin.RouterGroup) {
	rg.POST("/albums", controllers.CreateAlbum)
	// Add other routes here
}
