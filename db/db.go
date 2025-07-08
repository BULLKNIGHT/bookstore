package db

import (
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const dbName = "bookstore"
const collectionName = "books"

var Collection *mongo.Collection

func Init() {
	dbURL := os.Getenv("DATABASE_URL")
	// client options
	optionClient := options.Client().ApplyURI(dbURL)

	// connect to mongoDB
	client, err := mongo.Connect(optionClient)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("MongoDB connected successfully")

	// collection instance
	Collection = client.Database(dbName).Collection(collectionName)

	fmt.Println("Collection instance is ready")
}
