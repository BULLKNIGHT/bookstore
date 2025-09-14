package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/BULLKNIGHT/bookstore/db"
	"github.com/BULLKNIGHT/bookstore/logger"
	"github.com/BULLKNIGHT/bookstore/models"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func getAllBooks() ([]models.Book, error) {
	filter := bson.M{}
	cursor, err := db.Collection.Find(context.Background(), filter)

	var books []models.Book

	if err != nil {
		return books, err
	}

	for cursor.Next(context.Background()) {
		var book models.Book

		if err = cursor.Decode(&book); err != nil {
			return books, err
		}

		books = append(books, book)
	}

	defer cursor.Close(context.Background())

	logger.Log.Info("All books fetched successfully!! ✅")

	return books, nil
}

func insertBook(book models.Book) (*mongo.InsertOneResult, error) {
	result, err := db.Collection.InsertOne(context.Background(), book)

	if err != nil {
		return result, err
	}

	logger.Log.WithField("id", result.InsertedID).Info("Book inserted successfully!! 👌")
	return result, nil
}

func updateBook(book models.Book) (*mongo.UpdateResult, error) {
	filter := bson.M{"_id": book.ID}
	update := bson.M{"$set": book}

	result, err := db.Collection.UpdateOne(context.Background(), filter, update)

	if err != nil {
		return result, err
	}

	logger.Log.WithField("modified_count", result.ModifiedCount).Info("Book updated successfully!! 👌")
	return result, nil
}

func deleteBook(bookId bson.ObjectID) (*mongo.DeleteResult, error) {
	filter := bson.M{"_id": bookId}
	result, err := db.Collection.DeleteOne(context.Background(), filter)

	if err != nil {
		return result, nil
	}

	logger.Log.WithField("delete_count", result.DeletedCount).Info("Book deleted successfully!! ✅")
	return result, nil
}

func deleteAllBooks() (*mongo.DeleteResult, error) {
	filter := bson.M{}
	result, err := db.Collection.DeleteMany(context.Background(), filter)

	if err != nil {
		return result, err
	}

	logger.Log.WithField("delete_count", result.DeletedCount).Info("All books deleted successfully!! ✅")
	return result, nil
}

func validateBook(r *http.Request) (models.Book, error) {
	// no json data send
	if r.Body == nil {
		return models.Book{}, errors.New("no data found")
	}

	var book models.Book
	err := json.NewDecoder(r.Body).Decode(&book)

	// error during parsing json data
	if err != nil {
		return models.Book{}, errors.New("invalid data")
	}

	// validate required field
	if !book.IsValid() {
		return models.Book{}, errors.New("all fields (title, author, price) are required")
	}

	return book, nil
}

func GetAllBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	books, err := getAllBooks()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	json.NewEncoder(w).Encode(books)
}

func CreateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	book, err := validateBook(r)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	book.ID = bson.NewObjectID()
	_, err = insertBook(book)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	json.NewEncoder(w).Encode(book)
}

func UpdateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	bookId, err := bson.ObjectIDFromHex(params["id"])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Invalid object id")
		return
	}

	book, err := validateBook(r)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	book.ID = bookId
	result, err := updateBook(book)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Log.WithError(err).Error(err.Error())
		json.NewEncoder(w).Encode(err)
		return
	}

	if result.MatchedCount == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode("no data found by given id")
		return
	}

	json.NewEncoder(w).Encode(book)
}

func DeleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	bookId, err := bson.ObjectIDFromHex(params["id"])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Invalid object id")
		return
	}

	result, err := deleteBook(bookId)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	if result.DeletedCount == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode("no data found by given id")
		return
	}

	json.NewEncoder(w).Encode("Data deleted successfully")
}

func DeleteAllBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	_, err := deleteAllBooks()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	json.NewEncoder(w).Encode("all books deleted successfully")
}

func ServeHome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to book API"))
}
