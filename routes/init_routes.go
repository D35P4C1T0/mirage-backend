package routes

import (
	"github.com/gin-gonic/gin"
	"mirage-backend/routes/other"
)

func InitRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")
	{
		SetupUserRoutes(api)
		SetupAlbumRoutes(api)
		SetupPictureRoutes(api)

		// Homepage result
		other.SetupHomepageRoutes(api)

		// Health check route
		other.SetupHealthRoutes(api)
	}
}
