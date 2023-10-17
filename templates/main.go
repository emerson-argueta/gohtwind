package main

import (
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not specified
	}

	http.Handle("/static/", LoggingMiddleware(http.StripPrefix("/static/", http.FileServer(http.Dir("./frontend/static/")))))
	// TODO: Setup routes and middleware

	log.Printf("Server started on :%s\n", port)
	http.ListenAndServe(":"+port, nil)
}
