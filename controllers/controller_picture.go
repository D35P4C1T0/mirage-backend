package controllers

import (
	"context"
	"errors"
	"io"
	"log"
	"mirage-backend/utils"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"mirage-backend/database"
	"mirage-backend/models"
)

const timeoutDuration = 5 * time.Second

const CompressionQuality = 80

const PictureCollectionName = "pictures"

var pictureCollection *mongo.Collection

// InitializePictureController initializes the picture collection
func InitializePictureController() {
	pictureCollection = database.GetCollection(PictureCollectionName)
}

// UploadPicture handles the upload of a new picture
func UploadPicture(c *gin.Context) {
	// Set a context timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
	defer cancel()

	// Ensure the correct content type
	contentType := c.GetHeader("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		c.JSON(http.StatusUnsupportedMediaType, gin.H{"error": "Content-Type must be an image"})
		return
	}

	// Read the raw body
	fileBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file content", "details": err.Error()})
		return
	}

	// Compress the image
	compressedImage, err := utils.ScaleAndConvertToWebPBytes(fileBytes, CompressionQuality)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to compress image", "details": err.Error()})
		return
	}

	// Create and save the picture object (or handle as needed)
	picture := models.Picture{
		ID:         primitive.NewObjectID(),
		ImageData:  compressedImage,
		UploadedAt: time.Now(),
	}

	if _, err := pictureCollection.InsertOne(ctx, picture); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save picture", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Picture uploaded successfully",
		"data":    picture,
	})
}

// GetPicturesInAlbum retrieves all pictures from an album
func GetPicturesInAlbum(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	albumID := c.Param("albumId")
	albumObjectID, err := primitive.ObjectIDFromHex(albumID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid album ID"})
		return
	}

	filter := bson.M{"album_id": albumObjectID}
	cursor, err := pictureCollection.Find(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve pictures", "details": err.Error()})
		return
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}(cursor, ctx)

	var pictures []models.Picture
	if err := cursor.All(ctx, &pictures); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode pictures", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pictures retrieved successfully", "data": pictures})
}

// GetPictureByID retrieves a specific picture by ID
func GetPictureByID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pictureID := c.Param("pictureId")
	pictureObjectID, err := primitive.ObjectIDFromHex(pictureID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid picture ID"})
		return
	}

	filter := bson.M{"_id": pictureObjectID}
	var picture models.Picture
	if err := pictureCollection.FindOne(ctx, filter).Decode(&picture); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pictureID := c.Param("pictureId")
	pictureObjectID, err := primitive.ObjectIDFromHex(pictureID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid picture ID"})
		return
	}

	filter := bson.M{"_id": pictureObjectID}
	result, err := pictureCollection.DeleteOne(ctx, filter)
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

// GetAllPictures retrieves all pictures
func GetAllPictures(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := pictureCollection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve pictures", "details": err.Error()})
		return
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			log.Println(err)
		}
	}(cursor, ctx)

	var pictures []models.Picture
	if err := cursor.All(ctx, &pictures); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode pictures", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pictures retrieved successfully", "data": pictures})
}

func DeCouplePictureFromAlbum(c *gin.Context) {
	albumID := c.Param("albumId")
	pictureID := c.Param("pictureId")

	// Validate album ID
	albumObjectID, err := primitive.ObjectIDFromHex(albumID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid album ID"})
		return
	}

	// Validate picture ID
	pictureObjectID, err := primitive.ObjectIDFromHex(pictureID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid picture ID"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Update the album to remove the picture ID from the PictureIDs list
	filter := bson.M{"_id": albumObjectID}
	update := bson.M{"$pull": bson.M{"pictureIds": pictureObjectID}}

	result, err := albumCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to dissociate picture from album", "details": err.Error()})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Album not found"})
		return
	}

	if result.ModifiedCount == 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Picture was not associated with the album"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Picture dissociated from album successfully"})
}
