package main

import (
	"balancer-api/handlers"
	"balancer-api/models"
	"balancer-api/services"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	environment = os.Getenv("GO_ENV")
	port        = os.Getenv("PORT")
	uiOrigin    = os.Getenv("BALANCER_UI_ORIGIN")
	dbURL       = os.Getenv("DATABASE_URL")
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

	// Set up DB
	var db *gorm.DB
	var err error
	if environment == "dev" {
		db, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
		if err != nil {
			panic("gorm failed to connect to the sqlite database")
		}
	} else {
		db, err = gorm.Open(mysql.Open(dbURL), &gorm.Config{})
		if err != nil {
			panic("gorm failed to connect to the mysql database")
		}
	}

	db.AutoMigrate(&models.Record{})

	// Create services
	rs := &services.RecordService{DB: db}

	// Attach services to handler
	var h handlers.Handler
	h.RecordService = rs

	// Define routes
	router.Route("/records", func(r chi.Router) {
		r.Get("/", h.GetAllRecords)
		r.Post("/", h.CreateRecord)
		r.Put("/{id}", h.UpdateRecord)
		r.Delete("/{id}", h.DeleteRecord)
		r.Get("/net", h.GetNetWorth)
		r.Get("/sum", h.GetTypeSum)
	})

	// Start server
	if port == "" {
		port = "8080"
	}
	log.Printf("Starting server on %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
