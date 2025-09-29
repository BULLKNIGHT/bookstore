package main

import (
	"net/http"

	"github.com/BULLKNIGHT/bookstore/db"
	"github.com/BULLKNIGHT/bookstore/logger"
	"github.com/BULLKNIGHT/bookstore/middlewares"
	"github.com/BULLKNIGHT/bookstore/otel"
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
}

func main() {
	_, err := otel.Init()

	if err != nil {
		logger.Log.WithError(err).Error("Failed to create tracer provider!! 👎")
		return
	}

	defer func() {
		otel.ShutDown()
	}()

	_, err = db.Init()

	if err != nil {
		logger.Log.WithError(err).Error("MongoDB connection failed!! 👎")
		return
	}

	defer func() {
		db.Disconnect()
	}()

	r := mux.NewRouter()

	r.Use(middlewares.RecoverMiddleware)
	r.Use(middlewares.TraceMiddleware)
	r.Use(middlewares.LoggerMiddleware)

	routes.RegisterBook(r)

	logger.Log.Fatal(http.ListenAndServe(":4000", r))
}
