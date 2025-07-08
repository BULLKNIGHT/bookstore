package routes

import (
	"github.com/BULLKNIGHT/bookstore/controllers"
	"github.com/gorilla/mux"
)

func RegisterBook(router *mux.Router) {
	router.HandleFunc("/", controllers.ServeHome).Methods("GET")
	router.HandleFunc("/getAllBooks", controllers.GetAllBooks).Methods("GET")
	router.HandleFunc("/createBook", controllers.CreateBook).Methods("POST")
	router.HandleFunc("/updateBook/{id}", controllers.UpdateBook).Methods("PUT")
	router.HandleFunc("/deleteBook/{id}", controllers.DeleteBook).Methods("DELETE")
	router.HandleFunc("/deleteAllBooks", controllers.DeleteAllBooks).Methods("DELETE")
}
