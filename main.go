package main

import (
	"balancer-api/db"
	"balancer-api/handlers"
	"balancer-api/models"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
)

var (
	port = os.Getenv("PORT")
)

func main() {
	// Define router
	router := chi.NewRouter()

	db.Setup()

	// TODO needed?
	db.DB.AutoMigrate(&models.Record{})

	// Define routes
	router.Route("/records", func(r chi.Router) {
		r.Get("/", handlers.GetAllRecords)
		r.Post("/", handlers.CreateRecord)
		r.Put("/{id}", handlers.UpdateRecord)
		r.Delete("/{id}", handlers.DeleteRecord)
	})

	// Start server
	if port == "" {
		port = "8080"
	}
	log.Printf("Starting server on %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
