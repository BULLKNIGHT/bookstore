package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Book struct {
	ID            primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title         string        `json:"title" bson:"title"`
	Author        string        `json:"author" bson:"author"`
	Isbn          string        `json:"isbn" bson:"isbn"`
	PublishedYear int           `json:"published_year" bson:"published_year"`
	Price         int           `json:"price" bson:"price"`
	Category      string        `json:"category" bson:"category"`
}

func (book *Book) IsValid() bool {
	return book.Title != "" && book.Author != "" && book.Price > 0
}
