package routes

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// SetupHomepageRoutes is the handler for the homepage route
func SetupHomepageRoutes(rg *gin.RouterGroup) {
	rg.GET("/", Homepage)
}

func Homepage(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message":   "Welcome to Mirage!",
		"status":    "online",
		"timestamp": time.Now().Unix(),
		"info": gin.H{
			"version":     "0.1.0",
			"environment": "development",
			"description": "A powerful and elegant web application framework",
		},
		"endpoints": []string{
			"/api/v1/users",
			"/api/v1/pictures",
			"/api/v1/albums",
			"/health",
		},
		"quickLinks": gin.H{
			"documentation": "/docs",
			"contact":       "/support",
			"github":        "https://github.com/D35P4C1T0/mirage-backend",
		},
	})
}
