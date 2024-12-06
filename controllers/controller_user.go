package controllers

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"mirage-backend/database"
	"mirage-backend/models"
)

// TODO let chatgpt improve this code

const UserCollectionName = "users"

var userCollection *mongo.Collection

// InitializeUserController initializes the user collection
func InitializeUserController() {
	userCollection = database.GetCollection(UserCollectionName)
}

func GetAllUsers(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := userCollection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching users"})
		return
	}

	var users []models.User
	if err = cursor.All(ctx, &users); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding users"})
		return
	}

	c.JSON(http.StatusOK, users)
}

// GetUserProfile retrieves user profile information.
func GetUserProfile(c *gin.Context) {
	userId := c.Param("userId")
	objID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var user models.User
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = userCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// TODO make it granular and every param optional

// UpdateUserProfile updates user profile information.
func UpdateUserProfile(c *gin.Context) {
	userId := c.Param("userId")
	objID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var updatedUser models.User
	if err := c.ShouldBindJSON(&updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"username":           updatedUser.Username,
			"email":              updatedUser.Email,
			"password":           updatedUser.Password,
			"user_profile_id":    updatedUser.UserProfileID,
			"profile_picture_id": updatedUser.ProfilePictureID,
			"albums_id":          updatedUser.AlbumsID,
			"updated_at":         time.Now(),
		},
	}

	result, err := userCollection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating user"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

// TODO: make sure it is cascading

// DeleteUser removes a user account.
func DeleteUser(c *gin.Context) {
	userId := c.Param("userId")
	objID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := userCollection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting user"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
