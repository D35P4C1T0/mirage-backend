package controllers

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"mirage-backend/controllers/dbutils"
	"mirage-backend/database"
	"mirage-backend/models"
	"mirage-backend/utils"
	"net/http"
	"time"
)

// TODO: add picture compression

// UploadProfilePicture godoc
// @Summary Upload a profile picture
// @Description Uploads a profile picture to the database
// @Tags profile, pictures, pfp
// @Accept multipart/form-data
// @Produce json
// @Param userId path string true "User ID"
// @Param file formData file true "Profile picture file"
// @Success 201 {object} models.ProfilePicture
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /profilepictures/user/{userId} [post]
func UploadProfilePicture(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
	// Set context with timeout
	defer cancel()

	// Ensure multipart/form-data is used
	contentType := c.ContentType()
	if contentType != "multipart/form-data" {
		c.JSON(http.StatusUnsupportedMediaType, gin.H{"error": "Content-Type must be multipart/form-data"})
		return
	}

	var userObjectID primitive.ObjectID
	userId := c.Param("userId")
	if userId != "" {
		var err error
		userObjectID, err = primitive.ObjectIDFromHex(userId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}
	}

	userExists, err := dbutils.CheckIfItemExists(c, database.UserCollection, userObjectID)
	if !userExists {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check if user exists"})
		return
	}

	// Get the file from the request
	fileBytes, readErr := utils.RetrieveImageFromHTTPForm(c, "file")
	if readErr {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file"})
		return
	}

	// Resize and compress the image
	// I don't know how it behaves if the
	// picture has a 1:1 aspect ratio - idk
	compressedImage, compressErr := utils.ScaleAndConvertToWebPBytes(fileBytes, 45)
	if compressErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process image", "details": compressErr.Error()})
		return
	}

	newPicture := models.Picture{
		ID:          primitive.NewObjectID(),
		Description: "Profile picture of user " + userId,
	}

	// Insert the picture into the database
	pictureAdded, err :=
		dbutils.AddPictureToDB(ctx, compressedImage, database.PictureDataCollection, database.PictureCollection, newPicture)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error uploading profile picture to the database"})
		return
	}
	if !pictureAdded {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload profile picture to the database"})
		return
	}

	// Create a new profile picture object
	profilePicture := models.ProfilePicture{
		PictureID: newPicture.ID,
		UserID:    userObjectID,
		CreatedAt: time.Now(),
	}

	// Insert the profile picture into the database
	result, err := database.PfpCollection.InsertOne(ctx, profilePicture)
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
