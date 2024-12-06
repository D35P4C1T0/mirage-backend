package database

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"mirage-backend/config"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Connection struct {
	Database *mongo.Database
	Client   *mongo.Client
}

var Db = Connection{}

func SetupDatabase() {
	mongoURI := config.GetDatabaseURI()
	if mongoURI == "" {
		log.Fatal("MONGO_URI not set in .env file")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	clientOptions := options.Client().ApplyURI(mongoURI)

	// Connect to MongoDB
	Db.Client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// set database
	Db.Database = Db.Client.Database(config.GetDatabaseName())

	// verify connection
	err = Db.Client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MongoDB!")
}

func GetCollection(collectionName string) *mongo.Collection {
	// Ensure 'users' collection exists
	collection, err := EnsureCollection(Db.Database, collectionName)
	if err != nil {
		log.Fatalf("Error ensuring users collection: %v", err)
	}
	return collection
}

// EnsureCollection checks if a collection exists and creates it if it doesn't
func EnsureCollection(db *mongo.Database, collectionName string) (*mongo.Collection, error) {
	// Context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Check existing collections
	collections, err := db.ListCollectionNames(ctx, bson.D{})
	if err != nil {
		return nil, fmt.Errorf("failed to list collections: %v", err)
	}

	// Check if collection exists
	collectionExists := false
	for _, collection := range collections {
		if collection == collectionName {
			collectionExists = true
			break
		}
	}

	// If collection doesn't exist, create it
	if !collectionExists {
		// Optional: Create with specific options
		opts := options.CreateCollection().SetCapped(false)

		err := db.CreateCollection(ctx, collectionName, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to create collection %s: %v", collectionName, err)
		}

		log.Printf("Created new collection: %s", collectionName)
	}

	// Return the collection
	return db.Collection(collectionName), nil
}
