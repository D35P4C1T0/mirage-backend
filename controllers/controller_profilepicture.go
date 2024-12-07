package controllers

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"mime/multipart"
	"mirage-backend/database"
	"mirage-backend/models"
	"net/http"
	"time"
)

const PfpCollectionName = "profilepictures"

var pfpCollection *mongo.Collection

// InitializeProfilePictureController initializes the user collection
func InitializeProfilePictureController() {
	pfpCollection = database.GetCollection(PfpCollectionName)
}

// TODO: add picture compression

// UploadProfilePicture handles the uploading of a new profile picture
func UploadProfilePicture(c *gin.Context) {
	userId := c.Param("userId")
	objID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file"})
		return
	}

	fileHeader, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error opening file"})
		return
	}

	defer func(fileHeader multipart.File) {
		err := fileHeader.Close()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error closing file"})
		}
	}(fileHeader)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Read the file
	fileBytes := make([]byte, file.Size)
	if _, err := fileHeader.Read(fileBytes); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading file"})
		return
	}

	// Create a new profile picture
	profilePicture := models.ProfilePicture{
		ImageData: fileBytes,
		UserID:    objID,
		CreatedAt: time.Now(),
	}

	result, err := userCollection.InsertOne(ctx, profilePicture)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error uploading profile picture"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Profile picture uploaded successfully", "id": result.InsertedID})
}
