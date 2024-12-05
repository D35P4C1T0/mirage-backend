package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// SetupHomepageRoutes is the handler for the homepage route
func SetupHomepageRoutes(rg *gin.RouterGroup) {
	rg.GET("/", Homepage)
}

func Homepage(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Welcome to Mirage!",
	})
}
