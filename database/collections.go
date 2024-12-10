package database

import "go.mongodb.org/mongo-driver/mongo"

var (
	PfpCollection         *mongo.Collection
	AlbumCollection       *mongo.Collection
	PictureCollection     *mongo.Collection
	PictureDataCollection *mongo.Collection
	UserCollection        *mongo.Collection
)

// Collection names
const (
	PfpCollectionName         = "profilepictures"
	AlbumCollectionName       = "albums"
	PictureCollectionName     = "pictures"
	PictureDataCollectionName = "pictureData"
	UserCollectionName        = "users"
)

// InitializeCollections initializes all MongoDB collections used in the application
func InitializeCollections() {
	PfpCollection = GetCollection(PfpCollectionName)
	AlbumCollection = GetCollection(AlbumCollectionName)
	PictureCollection = GetCollection(PictureCollectionName)
	PictureDataCollection = GetCollection(PictureDataCollectionName)
	UserCollection = GetCollection(UserCollectionName)
}
