package main

import (
	"log"
	"net/http"
	"os"

	"github.com/BULLKNIGHT/bookstore/db"
	_ "github.com/BULLKNIGHT/bookstore/docs"
	"github.com/BULLKNIGHT/bookstore/logger"
	"github.com/BULLKNIGHT/bookstore/middlewares"
	"github.com/BULLKNIGHT/bookstore/otel"
	"github.com/BULLKNIGHT/bookstore/routes"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
)

// @title Bookstore API
// @version 1.0
// @description This is a Bookstore server with JWT authentication.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:4000
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
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
	r.Use(middlewares.RateLimiterMiddleware)
	r.Use(otelmux.Middleware("bookstore-api"))
	r.Use(middlewares.LoggerMiddleware)

	routes.RegisterBook(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
		logger.Log.Info("Using default port 4000")
	}

	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		logger.Log.WithError(err).Error("Server failed to start")
	}
}
