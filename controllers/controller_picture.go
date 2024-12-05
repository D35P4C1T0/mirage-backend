package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"mirage-backend/database"
	"mirage-backend/models"
)

const DbName = "mirage"

// UploadPicture handles the uploading of a new picture to an album
func UploadPicture(c *gin.Context) {
	albumID := c.Param("albumId")
	var picture models.Picture

	// Validate JSON input
	if err := c.ShouldBindJSON(&picture); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	// Validate and convert album ID
	albumObjectID, err := primitive.ObjectIDFromHex(albumID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid album ID"})
		return
	}

	// Set picture properties
	picture.ID = primitive.NewObjectID()
	picture.AlbumID = albumObjectID
	picture.UploadedAt = time.Now()

	collection := database.Db.Database(DbName).Collection("pictures")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Insert picture into the database
	_, err = collection.InsertOne(ctx, picture)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload picture", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Picture uploaded successfully", "data": picture})
}

// GetPictures retrieves all pictures from an album
func GetPictures(c *gin.Context) {
	albumID := c.Param("albumId")
	collection := database.Db.Database(DbName).Collection("pictures")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Validate and convert album ID
	albumObjectID, err := primitive.ObjectIDFromHex(albumID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid album ID"})
		return
	}

	// Query the pictures collection
	filter := bson.M{"album_id": albumObjectID}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve pictures", "details": err.Error()})
		return
	}
	defer cursor.Close(ctx)

	// Collect pictures into a slice
	var pictures []models.Picture
	for cursor.Next(ctx) {
		var picture models.Picture
		if err := cursor.Decode(&picture); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode picture", "details": err.Error()})
			return
		}
		pictures = append(pictures, picture)
	}

	if err := cursor.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cursor error", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pictures retrieved successfully", "data": pictures})
}

// GetPicture retrieves a specific picture by ID
func GetPicture(c *gin.Context) {
	pictureID := c.Param("pictureId")
	collection := database.Db.Database(DbName).Collection("pictures")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Validate and convert picture ID
	pictureObjectID, err := primitive.ObjectIDFromHex(pictureID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid picture ID"})
		return
	}

	// Query the pictures collection
	filter := bson.M{"_id": pictureObjectID}
	var picture models.Picture
	err = collection.FindOne(ctx, filter).Decode(&picture)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Picture not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve picture", "details": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Picture retrieved successfully", "data": picture})
}

// DeletePicture handles the deletion of a specific picture by ID
func DeletePicture(c *gin.Context) {
	pictureID := c.Param("pictureId")
	collection := database.Db.Database(DbName).Collection("pictures")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Validate and convert picture ID
	pictureObjectID, err := primitive.ObjectIDFromHex(pictureID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid picture ID"})
		return
	}

	// Delete the picture
	filter := bson.M{"_id": pictureObjectID}
	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete picture", "details": err.Error()})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Picture not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Picture deleted successfully"})
}
