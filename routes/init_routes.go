package routes

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"mirage-backend/routes/other"
	"mirage-backend/utils"
)

func InitRoutes(r *gin.Engine, apiPath string) {
	api := r.Group(apiPath)
	{
		SetupUserRoutes(api)
		SetupAlbumRoutes(api)
		SetupPictureRoutes(api)
		SetupProfilePictureRoutes(api)

		// Homepage result
		other.SetupHomepageRoutes(api)

		// Health check route
		other.SetupHealthRoutes(api)

		// Swagger documentation route

		utils.SetupDocs()
		// use ginSwagger middleware to serve the API docs
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
}
