package routes

import (
	"net/http"

	"github.com/BULLKNIGHT/bookstore/controllers"
	"github.com/BULLKNIGHT/bookstore/middlewares"
	"github.com/gorilla/mux"
)

func RegisterBook(router *mux.Router) {
	router.HandleFunc("/", controllers.ServeHome).Methods("GET")
	router.HandleFunc("/generateToken", controllers.GenerateToken).Methods("POST")
	router.Handle("/getAllBooks", middlewares.Chain(
		http.HandlerFunc(controllers.GetAllBooks),
		middlewares.AuthMiddleware),
	).Methods("GET")
	router.Handle("/createBook", middlewares.Chain(
		http.HandlerFunc(controllers.CreateBook),
		middlewares.AuthMiddleware,
		middlewares.RoleMiddleware("admin")),
	).Methods("POST")
	router.Handle("/updateBook/{id}", middlewares.Chain(
		http.HandlerFunc(controllers.UpdateBook),
		middlewares.AuthMiddleware,
		middlewares.RoleMiddleware("admin")),
	).Methods("PUT")
	router.Handle("/deleteBook/{id}", middlewares.Chain(
		http.HandlerFunc(controllers.DeleteBook),
		middlewares.AuthMiddleware,
		middlewares.RoleMiddleware("admin")),
	).Methods("DELETE")
	router.Handle("/deleteAllBooks", middlewares.Chain(
		http.HandlerFunc(controllers.DeleteAllBooks),
		middlewares.AuthMiddleware,
		middlewares.RoleMiddleware("admin")),
	).Methods("DELETE")
}
