package routes

import (
	"net/http"

	"github.com/BULLKNIGHT/bookstore/controllers"
	"github.com/BULLKNIGHT/bookstore/middlewares"
	"github.com/gorilla/mux"
)

func RegisterBook(router *mux.Router) {
	router.HandleFunc("/health", controllers.ServeHome).Methods("GET")
	router.HandleFunc("/token", controllers.GenerateToken).Methods("POST")
	router.Handle("/books", middlewares.Chain(
		http.HandlerFunc(controllers.GetAllBooks),
		middlewares.AuthMiddleware),
	).Methods("GET")
	router.Handle("/book", middlewares.Chain(
		http.HandlerFunc(controllers.CreateBook),
		middlewares.AuthMiddleware,
		middlewares.RoleMiddleware("admin")),
	).Methods("POST")
	router.Handle("/book/{id}", middlewares.Chain(
		http.HandlerFunc(controllers.UpdateBook),
		middlewares.AuthMiddleware,
		middlewares.RoleMiddleware("admin")),
	).Methods("PUT")
	router.Handle("/book/{id}", middlewares.Chain(
		http.HandlerFunc(controllers.DeleteBook),
		middlewares.AuthMiddleware,
		middlewares.RoleMiddleware("admin")),
	).Methods("DELETE")
	router.Handle("/books", middlewares.Chain(
		http.HandlerFunc(controllers.DeleteAllBooks),
		middlewares.AuthMiddleware,
		middlewares.RoleMiddleware("admin")),
	).Methods("DELETE")
}
