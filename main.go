package main

import (
	"net/http"

	"github.com/BULLKNIGHT/bookstore/db"
	"github.com/BULLKNIGHT/bookstore/logger"
	"github.com/BULLKNIGHT/bookstore/middlewares"
	"github.com/BULLKNIGHT/bookstore/routes"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		logger.Log.WithError(err).Error("No .env file found (using system env)")
	}

	logger.Init()
	db.Init()
}

func main() {
	r := mux.NewRouter()

	r.Use(middlewares.RecoverMiddleware)
	r.Use(middlewares.LoggerMiddleware)
	routes.RegisterBook(r)

	logger.Log.Fatal(http.ListenAndServe(":4000", r))
}
