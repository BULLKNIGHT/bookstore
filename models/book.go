package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Book represents a book in the bookstore
// @Description Book information with details like title, author, price, etc.
type Book struct {
	ID            primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty" swaggerignore:"true"`
	Title         string             `json:"title" bson:"title" example:"The Go Programming Language"`
	Author        string             `json:"author" bson:"author" example:"Charles Babage"`
	Isbn          string             `json:"isbn" bson:"isbn" example:"978-0134190440"`
	PublishedYear int                `json:"published_year" bson:"published_year" example:"2015"`
	Price         int                `json:"price" bson:"price" example:"2999"`
	Category      string             `json:"category" bson:"category" example:"Programming"`
}

func (book *Book) IsValid() bool {
	return book.Title != "" && book.Author != "" && book.Price > 0
}
