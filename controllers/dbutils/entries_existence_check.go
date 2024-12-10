package dbutils

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// CheckIfItemExists checks if an item exists in a collection
func CheckIfItemExists(ctx context.Context, collection *mongo.Collection, ID primitive.ObjectID) (bool, error) {
	result, err := collection.Find(ctx, primitive.M{"_id": ID})
	if err != nil {
		return false, err
	}
	if result.RemainingBatchLength() == 0 {
		return false, nil
	}
	return true, nil
}
