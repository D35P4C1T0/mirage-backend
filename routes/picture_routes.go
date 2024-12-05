package routes

import (
	"github.com/gin-gonic/gin"
	"mirage-backend/controllers"
)

// SetupPictureRoutes sets up all routes related to Picture
func SetupPictureRoutes(rg *gin.RouterGroup) {
	rg.GET("/pictures", controllers.GetPictures)
	rg.GET("/pictures/:pictureId", controllers.GetPicture)
	rg.DELETE("/pictures/:pictureId", controllers.DeletePicture)
	rg.POST("/pictures", controllers.UploadPicture)
	// Add other routes here
}
