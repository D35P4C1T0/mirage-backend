package routes

import (
	"github.com/gin-gonic/gin"
)

func InitRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		SetupHomepageRoutes(api)
		SetupUserRoutes(api)
		SetupAlbumRoutes(api)
		SetupPictureRoutes(api)
	}
}
