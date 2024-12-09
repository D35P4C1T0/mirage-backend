package controllers

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"mirage-backend/database"
	"mirage-backend/models"
)

const AlbumCollectionName = "albums"

var albumCollection *mongo.Collection

// InitializeAlbumController initializes the album collection
func InitializeAlbumController() {
	albumCollection = database.GetCollection(AlbumCollectionName)
}

// CreateAlbum godoc
// @Summary Create a new album
// @Description Creates a new album with the provided details
// @Tags albums
// @Accept json
// @Produce json
// @Param album body models.Album true "Album to create"
// @Success 201 {object} map[string]interface{} "Album created successfully"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Failed to create album"
// @Router /albums [post]
func CreateAlbum(c *gin.Context) {
	var album models.Album

	// Validate JSON input
	if err := c.ShouldBindJSON(&album); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	// Set album properties
	album.ID = primitive.NewObjectID()
	album.CreatedAt = time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Insert album into the database
	_, err := albumCollection.InsertOne(ctx, album)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create album", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Album created successfully", "data": album})
}

// GetAllAlbums godoc
// @Summary Retrieve all albums
// @Description Fetches a list of all albums in the database
// @Tags albums
// @Produce json
// @Success 200 {object} map[string]interface{} "Albums retrieved successfully"
// @Failure 500 {object} map[string]string "Failed to retrieve albums"
// @Router /albums [get]
func GetAllAlbums(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := albumCollection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve albums", "details": err.Error()})
		return
	}

	defer func(cursor *mongo.Cursor, ctx context.Context) {
		if err := cursor.Close(ctx); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to close cursor", "details": err.Error()})
		}
	}(cursor, ctx)

	// Collect albums into a slice
	var albums []models.Album
	for cursor.Next(ctx) {
		var album models.Album
		if err := cursor.Decode(&album); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode album", "details": err.Error()})
			return
		}
		albums = append(albums, album)
	}

	if err := cursor.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cursor error", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Albums retrieved successfully", "data": albums})
}

// GetAlbumByID godoc
// @Summary Retrieve a specific album
// @Description Fetches a single album by its unique identifier
// @Tags albums
// @Produce json
// @Param albumId path string true "Album Unique Identifier"
// @Success 200 {object} map[string]interface{} "Album retrieved successfully"
// @Failure 400 {object} map[string]string "Invalid album ID"
// @Failure 404 {object} map[string]string "Album not found"
// @Failure 500 {object} map[string]string "Failed to retrieve album"
// @Router /albums/{albumId} [get]
func GetAlbumByID(c *gin.Context) {
	albumID := c.Param("albumId")

	// Validate and convert album ID
	albumObjectID, err := primitive.ObjectIDFromHex(albumID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid album ID"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Query the albums collection
	var album models.Album
	filter := bson.M{"_id": albumObjectID}
	err = albumCollection.FindOne(ctx, filter).Decode(&album)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Album not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve album", "details": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Album retrieved successfully", "data": album})
}

// GetAlbumsByUserID godoc
// @Summary Retrieve albums for a specific user
// @Description Fetches all albums associated with a given user ID
// @Tags albums
// @Produce json
// @Param userId path string true "User Unique Identifier"
// @Success 200 {object} map[string]interface{} "Albums retrieved successfully"
// @Failure 400 {object} map[string]string "Invalid user ID"
// @Failure 500 {object} map[string]string "Failed to retrieve albums"
// @Router /albums/user/{userId} [get]
func GetAlbumsByUserID(c *gin.Context) {
	userID := c.Param("userId")

	// Validate and convert user ID if necessary (e.g., ensure it's a valid ObjectID)
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Query albums collection to find albums by user ID
	filter := bson.M{"user_id": userObjectID}
	cursor, err := albumCollection.Find(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve albums", "details": err.Error()})
		return
	}

	defer func(cursor *mongo.Cursor, ctx context.Context) {
		if err := cursor.Close(ctx); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to close cursor", "details": err.Error()})
		}
	}(cursor, ctx)

	// Collect albums into a slice
	var albums []models.Album
	for cursor.Next(ctx) {
		var album models.Album
		if err := cursor.Decode(&album); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode album", "details": err.Error()})
			return
		}
		albums = append(albums, album)
	}

	if err := cursor.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cursor error", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Albums retrieved successfully", "data": albums})
}

// UpdateAlbum godoc
// @Summary Update an existing album
// @Description Updates the details of a specific album by its ID
// @Tags albums
// @Accept json
// @Produce json
// @Param albumId path string true "Album Unique Identifier"
// @Param album body models.Album true "Album update information"
// @Success 200 {object} map[string]string "Album updated successfully"
// @Failure 400 {object} map[string]string "Invalid input or album ID"
// @Failure 404 {object} map[string]string "Album not found"
// @Failure 500 {object} map[string]string "Failed to update album"
// @Router /albums/{albumId} [put]
func UpdateAlbum(c *gin.Context) {
	albumID := c.Param("albumId")
	var updatedAlbum models.Album

	// Bind JSON payload to the album struct
	if err := c.ShouldBindJSON(&updatedAlbum); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	// Validate album ID
	albumObjectID, err := primitive.ObjectIDFromHex(albumID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid album ID"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Update album in the database
	filter := bson.M{"_id": albumObjectID}
	update := bson.M{"$set": updatedAlbum}
	result, err := albumCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update album", "details": err.Error()})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Album not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Album updated successfully"})
}

// DeleteAlbum godoc
// @Summary Delete an album
// @Description Permanently removes an album from the database by its ID
// @Tags albums
// @Produce json
// @Param albumId path string true "Album Unique Identifier"
// @Success 200 {object} map[string]string "Album deleted successfully"
// @Failure 400 {object} map[string]string "Invalid album ID"
// @Failure 404 {object} map[string]string "Album not found"
// @Failure 500 {object} map[string]string "Failed to delete album"
// @Router /albums/{albumId} [delete]
func DeleteAlbum(c *gin.Context) {
	albumID := c.Param("albumId")

	// Validate album ID
	albumObjectID, err := primitive.ObjectIDFromHex(albumID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid album ID"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Delete album
	filter := bson.M{"_id": albumObjectID}
	result, err := albumCollection.DeleteOne(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete album", "details": err.Error()})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Album not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Album deleted successfully"})
}

// SearchAlbums godoc
// @Summary Search albums
// @Description Searches albums by title using case-insensitive partial matching
// @Tags albums
// @Produce json
// @Param q query string true "Search query"
// @Success 200 {object} map[string]interface{} "Albums retrieved successfully"
// @Failure 500 {object} map[string]string "Failed to search albums"
// @Router /albums/search [get]
func SearchAlbums(c *gin.Context) {
	query := c.Query("q")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Search query using case-insensitive regex
	filter := bson.M{"title": bson.M{"$regex": query, "$options": "i"}}
	cursor, err := albumCollection.Find(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search albums", "details": err.Error()})
		return
	}

	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to close cursor", "details": err.Error()})
		}
	}(cursor, ctx)

	// Collect albums into a slice
	var albums []models.Album
	for cursor.Next(ctx) {
		var album models.Album
		if err := cursor.Decode(&album); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode album", "details": err.Error()})
			return
		}
		albums = append(albums, album)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Albums retrieved successfully", "data": albums})
}
