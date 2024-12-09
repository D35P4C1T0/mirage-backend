package controllers

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"io"
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

// UploadProfilePicture godoc
// @Summary Upload a profile picture
// @Description Uploads a profile picture to the database
// @Tags profilepictures
// @Accept multipart/form-data
// @Produce json
// @Param userId path string true "User ID"
// @Param file formData file true "Profile picture file"
// @Success 201 {object} models.ProfilePicture
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /profilepictures/user/{userId} [post]
func UploadProfilePicture(c *gin.Context) {
	// Parse and validate user ID
	userId := c.Param("userId")
	objID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	// Get the uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File not provided or invalid"})
		return
	}

	// Open the uploaded file
	fileHeader, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open uploaded file"})
		return
	}
	defer func(fileHeader multipart.File) {
		err := fileHeader.Close()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to close uploaded file"})
		}
	}(fileHeader)

	// Read file content into a byte slice
	fileBytes, err := io.ReadAll(fileHeader)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file content"})
		return
	}

	// Context with timeout for database operations
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create a new profile picture object
	profilePicture := models.ProfilePicture{
		ImageData: fileBytes,
		UserID:    objID,
		CreatedAt: time.Now(),
	}

	// Insert the profile picture into the database
	result, err := userCollection.InsertOne(ctx, profilePicture)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload profile picture to the database"})
		return
	}

	// Return success response
	c.JSON(http.StatusCreated, gin.H{
		"message": "Profile picture uploaded successfully",
		"id":      result.InsertedID,
	})
}
