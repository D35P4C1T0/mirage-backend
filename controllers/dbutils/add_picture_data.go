package dbutils

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"mirage-backend/models"
	"mirage-backend/utils"
)

// AddPictureToDB inserts a picture and its associated data into MongoDB collections.
//
// Parameters:
//   - c: The context for database operations.
//   - data: The byte slice containing the picture data.
//   - pictureDataCollection: The MongoDB collection to store picture data.
//   - pictureCollection: The MongoDB collection to store picture metadata.
//   - picture: The picture metadata to be stored.
//
// Returns:
//   - A boolean indicating success or failure of the operation.
//   - An error if any issue occurs during database operations or processing.
//
// This function performs the following steps:
//  1. Retrieves the dimensions (width and height) of the picture using the `utils.GetPictureDimensions` function.
//  2. Stores the compressed image data in the `pictureDataCollection`.
//  3. Updates the `Picture` model with the generated `PictureDataID`, width, and height.
//  4. Stores the picture metadata in the `pictureCollection`.
//
// In case of any error during dimension retrieval or database insertion, the function logs the error
// and returns `false` along with the error.
func AddPictureToDB(
	c context.Context,
	data []byte,
	pictureDataCollection *mongo.Collection,
	pictureCollection *mongo.Collection,
	picture models.Picture,
) (bool, error) {

	width, height, dimErr := utils.GetPictureDimensions(data)
	if dimErr != nil {
		log.Fatal(dimErr)
		return false, dimErr
	}

	// Store the compressed image in the database
	pictureData := models.PictureData{
		ID:   primitive.NewObjectID(),
		Data: data,
	}

	// add picture data to db
	if _, dbErr := pictureDataCollection.InsertOne(c, pictureData); dbErr != nil {
		log.Fatal(dbErr)
		return false, dbErr
	}

	picture.PictureDataID = pictureData.ID
	picture.Height = height
	picture.Width = width

	// add picture data to db
	if _, dbErr := pictureCollection.InsertOne(c, picture); dbErr != nil {
		log.Fatal(dbErr)
		return false, dbErr
	}

	return true, nil
}
