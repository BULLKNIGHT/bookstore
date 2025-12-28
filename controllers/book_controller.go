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
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func getAllBooks(ctx context.Context) ([]models.Book, error) {
	filter := bson.M{}
	cursor, err := db.Collection.Find(ctx, filter)

	var books []models.Book

	if err != nil {
		return books, err
	}

	for cursor.Next(ctx) {
		var book models.Book

		if err = cursor.Decode(&book); err != nil {
			return books, err
		}

		books = append(books, book)
	}

	defer cursor.Close(ctx)

	logger.Log.Info("All books fetched successfully!! ✅")

	return books, nil
}

func insertBook(book models.Book, ctx context.Context) (*mongo.InsertOneResult, error) {
	result, err := db.Collection.InsertOne(ctx, book)

	if err != nil {
		return result, err
	}

	logger.Log.WithField("id", result.InsertedID).Info("Book inserted successfully!! 👌")
	return result, nil
}

func updateBook(book models.Book, ctx context.Context) (*mongo.UpdateResult, error) {
	filter := bson.M{"_id": book.ID}
	update := bson.M{"$set": book}

	result, err := db.Collection.UpdateOne(ctx, filter, update)

	if err != nil {
		return result, err
	}

	logger.Log.WithField("modified_count", result.ModifiedCount).Info("Book updated successfully!! 👌")
	return result, nil
}
func deleteBook(bookId primitive.ObjectID, ctx context.Context) (*mongo.DeleteResult, error) {
	filter := bson.M{"_id": bookId}
	result, err := db.Collection.DeleteOne(ctx, filter)

	if err != nil {
		return result, err
	}

	logger.Log.WithField("delete_count", result.DeletedCount).Info("Book deleted successfully!! ✅")
	return result, nil
}

func deleteAllBooks(ctx context.Context) (*mongo.DeleteResult, error) {
	filter := bson.M{}
	result, err := db.Collection.DeleteMany(ctx, filter)

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

// GetAllBooks godoc
// @Summary Get all books
// @Description Retrieve all books from the database
// @Tags books
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.Book
// @Failure 401 {object} string "Unauthorized"
// @Failure 500 {object} string "Internal server error"
// @Router /books [get]
func GetAllBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	books, err := getAllBooks(r.Context())

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	json.NewEncoder(w).Encode(books)
}

// CreateBook godoc
// @Summary Create a new book
// @Description Add a new book to the database (Admin only)
// @Tags books
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param book body models.Book true "Book object"
// @Success 200 {object} models.Book
// @Failure 400 {object} string "Bad request"
// @Failure 401 {object} string "Unauthorized"
// @Failure 403 {object} string "Forbidden - Admin role required"
// @Failure 500 {object} string "Internal server error"
// @Router /book [post]
func CreateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	book, err := validateBook(r)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	book.ID = primitive.NewObjectID()
	_, err = insertBook(book, r.Context())

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	json.NewEncoder(w).Encode(book)
}

// UpdateBook godoc
// @Summary Update a book
// @Description Update an existing book by ID (Admin only)
// @Tags books
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Book ID"
// @Param book body models.Book true "Book object"
// @Success 200 {object} models.Book
// @Failure 400 {object} string "Bad request"
// @Failure 401 {object} string "Unauthorized"
// @Failure 403 {object} string "Forbidden - Admin role required"
// @Failure 404 {object} string "Book not found"
// @Failure 500 {object} string "Internal server error"
// @Router /book/{id} [put]
func UpdateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	bookId, err := primitive.ObjectIDFromHex(params["id"])

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
	result, err := updateBook(book, r.Context())

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

// DeleteBook godoc
// @Summary Delete a book
// @Description Delete a book by ID (Admin only)
// @Tags books
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Book ID"
// @Success 200 {object} string "Data deleted successfully"
// @Failure 400 {object} string "Bad request"
// @Failure 401 {object} string "Unauthorized"
// @Failure 403 {object} string "Forbidden - Admin role required"
// @Failure 404 {object} string "Book not found"
// @Failure 500 {object} string "Internal server error"
// @Router /book/{id} [delete]
func DeleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	bookId, err := primitive.ObjectIDFromHex(params["id"])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Invalid object id")
		return
	}

	result, err := deleteBook(bookId, r.Context())

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

// DeleteAllBooks godoc
// @Summary Delete all books
// @Description Delete all books from the database (Admin only)
// @Tags books
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} string "All books deleted successfully"
// @Failure 401 {object} string "Unauthorized"
// @Failure 403 {object} string "Forbidden - Admin role required"
// @Failure 500 {object} string "Internal server error"
// @Router /books [delete]
func DeleteAllBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	_, err := deleteAllBooks(r.Context())

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	json.NewEncoder(w).Encode("all books deleted successfully")
}

// ServeHome godoc
// @Summary Home page
// @Description Welcome message for the API
// @Tags general
// @Produce plain
// @Success 200 {string} string "Welcome to bookstore API"
// @Router /health [get]
func ServeHome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to bookstore API"))
}
