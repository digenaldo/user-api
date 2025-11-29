package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"

	httphandler "user-api/internal/handler/http"
	"user-api/internal/infra/mongo"
	"user-api/internal/repository"
	"user-api/internal/usecase"
)

func main() {
	// Read environment variables
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Create MongoDB client
	client := mongo.NewClient(mongoURI)
	defer func() {
		if err := client.Disconnect(nil); err != nil {
			log.Printf("Error disconnecting from MongoDB: %v", err)
		}
	}()

	// Get database
	db := client.Database("userdb")

	// Create repository
	repo := repository.NewUserMongoRepository(db)

	// Create use case
	uc := usecase.NewUserUseCase(repo)

	// Create handler
	handler := httphandler.NewUserHandler(uc)

	// Setup router
	r := chi.NewRouter()
	handler.RegisterRoutes(r)

	// Start server
	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
