package controllers

import (
	"bytes"
	"context"
	"errors"
	"image"
	"io"
	"log"
	"mime/multipart"
	"mirage-backend/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"mirage-backend/database"
	"mirage-backend/models"
)

const timeoutDuration = 10 * time.Second

const CompressionQuality = 80

const PictureCollectionName = "pictures"

var pictureCollection *mongo.Collection

// InitializePictureController initializes the picture collection
func InitializePictureController() {
	pictureCollection = database.GetCollection(PictureCollectionName)
}

// UploadPicture handles uploading a new picture via a multipart/form-data POST request.
func UploadPicture(c *gin.Context) {
	// Set context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
	defer cancel()

	// Ensure multipart/form-data is used
	contentType := c.ContentType()
	if contentType != "multipart/form-data" {
		c.JSON(http.StatusUnsupportedMediaType, gin.H{"error": "Content-Type must be multipart/form-data"})
		return
	}

	// Retrieve the uploaded file
	file, header, fileErr := c.Request.FormFile("file")
	if fileErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to retrieve file", "details": fileErr.Error()})
		return
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to close file", "details": err.Error()})
		}
	}(file)

	log.Printf("Received file: %s, size: %d bytes", header.Filename, header.Size)

	// Read the file into memory
	fileBytes, readErr := io.ReadAll(file)
	if readErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file", "details": readErr.Error()})
		return
	}

	// Decode and validate the image
	_, format, decodeErr := image.Decode(bytes.NewReader(fileBytes))
	if decodeErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid image format", "details": decodeErr.Error()})
		return
	}

	log.Printf("Decoded image format: %s", format)

	// Resize and compress the image
	compressedImage, compressErr := utils.ScaleAndConvertToWebPBytes(fileBytes, CompressionQuality)
	if compressErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process image", "details": compressErr.Error()})
		return
	}

	// Store the compressed image in the database
	picture := models.Picture{
		ID:         primitive.NewObjectID(),
		ImageData:  compressedImage,
		UploadedAt: time.Now(),
	}

	if _, dbErr := pictureCollection.InsertOne(ctx, picture); dbErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload picture", "details": dbErr.Error()})
		return
	}

	// Return success response
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
