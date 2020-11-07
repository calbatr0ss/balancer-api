package main

import (
	"balancer-api/db"
	"balancer-api/handlers"
	"balancer-api/models"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

var (
	port     = os.Getenv("PORT")
	uiOrigin = os.Getenv("BALANCER_UI_ORIGIN")
)

func main() {
	// Define router
	router := chi.NewRouter()
	// Allow cors from frontend
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{uiOrigin},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	db.Setup()

	// TODO needed?
	db.DB.AutoMigrate(&models.Record{})

	// Define routes
	router.Route("/records", func(r chi.Router) {
		r.Get("/", handlers.GetAllRecords)
		r.Post("/", handlers.CreateRecord)
		r.Put("/{id}", handlers.UpdateRecord)
		r.Delete("/{id}", handlers.DeleteRecord)
		r.Get("/net", handlers.GetNetWorth)
		r.Get("/sum", handlers.GetTypeSum)
	})

	// Start server
	if port == "" {
		port = "8080"
	}
	log.Printf("Starting server on %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
