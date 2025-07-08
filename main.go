package main

import (
	"log"
	"net/http"

	"github.com/BULLKNIGHT/bookstore/db"
	"github.com/BULLKNIGHT/bookstore/routes"
	"github.com/gorilla/mux"
)

func main() {
	db.Init()

	r := mux.NewRouter()
	routes.RegisterBook(r)

	log.Fatal(http.ListenAndServe(":4000", r))
}
