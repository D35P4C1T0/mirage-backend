package controllers

import (
	"context"
	"errors"
	"log"
	"mirage-backend/controllers/dbutils"
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

// UploadPicture godoc
// @Summary Upload a picture
// @Description Uploads a picture to the database
// @Tags pictures
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Picture file"
// @Param albumId path string false "Album ID"
// @Success 201 {object} models.Picture
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /pictures [post]
// @Router /albums/{albumId}/pictures [post]
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

	// Get the album ID from the request
	var albumObjectID primitive.ObjectID
	albumId := c.Param("albumId")
	if albumId != "" {
		var err error
		albumObjectID, err = primitive.ObjectIDFromHex(albumId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid album ID"})
			return
		}
	}

	// check if album exists
	albumExists, err := dbutils.CheckIfItemExists(c, database.AlbumCollection, albumObjectID)
	if !albumExists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Album not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check if album exists in the database"})
	}

	//// for the future
	//userExists, err := dbutils.CheckIfItemExists(c, userCollection, userObjectID)
	//if !userExists {
	//	c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
	//	return
	//} else if err != nil {
	//	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check if user exists"})
	//	return
	//}

	// Get the file from the request
	fileBytes, readErr := utils.RetrieveImageFromHTTPForm(c, "file")
	if readErr {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file"})
		return
	}

	//log.Printf("Decoded image format: %s", format)

	// Resize and compress the image
	compressedImage, compressErr := utils.ScaleAndConvertToWebPBytes(fileBytes, CompressionQuality)
	if compressErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process image", "details": compressErr.Error()})
		return
	}

	// Store the compressed image in the database
	picture := models.Picture{
		ID:         primitive.NewObjectID(),
		UploadedAt: time.Now(),
		AlbumID:    albumObjectID,
	}

	// Insert the picture into the database
	pictureAdded, err :=
		dbutils.AddPictureToDB(ctx, compressedImage, database.PictureDataCollection, database.PictureCollection, picture)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error uploading profile picture to the database"})
		return
	}
	if !pictureAdded {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload profile picture to the database"})
		return
	}

	// Return success response
	c.JSON(http.StatusCreated, gin.H{
		"message": "Picture uploaded successfully",
		"data":    picture,
	})
}

// GetPictureData godoc
// @Summary Get picture data
// @Description Retrieves the raw image data of a specific picture
// @Tags pictures
// @Accept json
// @Produce image/webp
// @Param pictureId path string true "Picture ID"
// @Success 200 {string} string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /pictures/{pictureId}/data [get]
func GetPictureData(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pictureID := c.Param("pictureId")
	pictureObjectID, err := primitive.ObjectIDFromHex(pictureID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid picture ID"})
		return
	}

	// retrieve the picture
	var picture models.Picture
	if err := database.PictureCollection.FindOne(ctx, bson.M{"_id": pictureObjectID}).Decode(&picture); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Picture not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve picture", "details": err.Error()})
		}
		return
	}

	// Retrieve the picture data
	var pictureData models.PictureData
	if err := database.PictureDataCollection.FindOne(ctx, bson.M{"_id": picture.PictureDataID}).Decode(&pictureData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve picture data", "details": err.Error()})
		return
	}

	// fetch the image data

	c.Data(http.StatusOK, "image/webp", pictureData.Data)
}

// GetPicturesInAlbum godoc
// @Summary Get pictures in an album
// @Description Retrieves all pictures in a specific album
// @Tags pictures
// @Accept json
// @Produce json
// @Param albumId path string true "Album ID"
// @Success 200 {object} []models.Picture
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /albums/{albumId}/pictures [get]
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
	cursor, err := database.PictureCollection.Find(ctx, filter)
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

// GetPictureByID godoc
// @Summary Get picture by ID
// @Description Retrieves a specific picture by its ID
// @Tags pictures
// @Accept json
// @Produce json
// @Param pictureId path string true "Picture ID"
// @Success 200 {object} models.Picture
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /pictures/{pictureId} [get]
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
	if err := database.PictureCollection.FindOne(ctx, filter).Decode(&picture); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Picture not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve picture", "details": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Picture retrieved successfully", "data": picture})
}

// DeletePicture godoc
// @Summary Delete picture by ID
// @Description Deletes a specific picture by its ID
// @Tags pictures
// @Accept json
// @Produce json
// @Param pictureId path string true "Picture ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /pictures/{pictureId} [delete]
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
	result, err := database.PictureCollection.DeleteOne(ctx, filter)
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

// GetAllPictures godoc
// @Summary Get all pictures
// @Description Retrieves all pictures from the database
// @Tags pictures
// @Accept json
// @Produce json
// @Success 200 {object} []models.Picture
// @Failure 500 {object} map[string]string
// @Router /pictures [get]
func GetAllPictures(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := database.PictureCollection.Find(ctx, bson.M{})
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

// DeCouplePictureFromAlbum godoc
// @Summary Remove picture from album
// @Description Removes a picture's association with a specific album without deleting the picture
// @Tags pictures
// @Accept json
// @Produce json
// @Param albumId path string true "Album ID"
// @Param pictureId path string true "Picture ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /albums/{albumId}/pictures/{pictureId} [delete]
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

	result, err := database.AlbumCollection.UpdateOne(ctx, filter, update)
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
