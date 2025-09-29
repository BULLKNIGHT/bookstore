package db

import (
	"context"
	"os"

	"github.com/BULLKNIGHT/bookstore/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo"
)

const dbName = "bookstore"
const collectionName = "books"

var Collection *mongo.Collection
var client *mongo.Client

func Init() (*mongo.Client, error) {
	dbURL := os.Getenv("MONGO_URL")
	// client options
	optionClient := options.Client().ApplyURI(dbURL).SetMonitor(otelmongo.NewMonitor())

	// connect to mongoDB
	client, err := mongo.Connect(context.Background(), optionClient)

	if err != nil {
		return nil, err
	}

	logger.Log.Info("MongoDB connected successfully!! 👍")

	// collection instance
	Collection = client.Database(dbName).Collection(collectionName)

	logger.Log.Info("Collection instance is ready!! 👌")

	return client, nil
}

func Disconnect() {
	if err := client.Disconnect(context.Background()); err != nil {
		logger.Log.WithError(err).Error("MongoDB failed to disconnect!! 👎")
	} else {
		logger.Log.Info("MongoDB disconnected gracefully!! 👍")
	}
}
