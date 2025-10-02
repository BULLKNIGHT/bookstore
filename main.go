package main

import (
	"log"
	"net/http"

	"github.com/BULLKNIGHT/bookstore/db"
	"github.com/BULLKNIGHT/bookstore/logger"
	"github.com/BULLKNIGHT/bookstore/middlewares"
	"github.com/BULLKNIGHT/bookstore/otel"
	"github.com/BULLKNIGHT/bookstore/routes"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// Load env variables
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found (using system env): %v", err)
	}

	// Initialize OpenTelemetry
	if err := otel.Init(); err != nil {
		log.Printf("OpenTelemetry failed to initiate: %v", err)
		return
	} else {
		defer func() {
			otel.ShutDown()
		}()
	}

	// Initialize logger
	logger.Init()

	// Initialize DB
	if _, err := db.Init(); err != nil {
		logger.Log.WithError(err).Error("MongoDB connection failed!! 👎")
		return
	} else {
		defer func() {
			db.Disconnect()
		}()
	}

	r := mux.NewRouter()

	r.Use(middlewares.RecoverMiddleware)
	r.Use(middlewares.TraceMiddleware)
	r.Use(middlewares.LoggerMiddleware)

	routes.RegisterBook(r)

	err := http.ListenAndServe(":4000", r)
	if err != nil {
		logger.Log.WithError(err).Error("Server failed to start")
	}
}
