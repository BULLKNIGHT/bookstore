package db

import (
	"os"

	"github.com/BULLKNIGHT/bookstore/logger"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const dbName = "bookstore"
const collectionName = "books"

var Collection *mongo.Collection

func Init() {
	dbURL := os.Getenv("MONGO_URL")
	// client options
	optionClient := options.Client().ApplyURI(dbURL)

	// connect to mongoDB
	client, err := mongo.Connect(optionClient)

	if err != nil {
		logger.Log.WithError(err).Fatal("MongoDB connection failed!! 👎")
	}

	logger.Log.Info("MongoDB connected successfully!! 👍")

	// collection instance
	Collection = client.Database(dbName).Collection(collectionName)

	logger.Log.Info("Collection instance is ready!! 👌")
}
