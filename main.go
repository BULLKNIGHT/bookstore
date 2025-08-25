package main

import (
	"log"
	"net/http"

	"github.com/BULLKNIGHT/bookstore/db"
	"github.com/BULLKNIGHT/bookstore/middlewares"
	"github.com/BULLKNIGHT/bookstore/routes"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found (using system env)")
	}

	db.Init()
}

func main() {
	r := mux.NewRouter()

	r.Use(middlewares.RecoverMiddleware)
	routes.RegisterBook(r)

	log.Fatal(http.ListenAndServe(":4000", r))
}
