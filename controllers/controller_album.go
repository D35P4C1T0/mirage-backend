package controllers

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"mirage-backend/database"
)

const AlbumCollectionName = "albums"

var albumCollection *mongo.Collection

// InitializeAlbumController initializes the picture collection
func InitializeAlbumController() {
	pictureCollection = database.GetCollection(AlbumCollectionName)
}

func CreateAlbum(c *gin.Context) {
	// TODO
}

func GetAllAlbums(c *gin.Context) {
	// TODO
}

func GetAlbumByID(c *gin.Context) {
	// TODO
}
