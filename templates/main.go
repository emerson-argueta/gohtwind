package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
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

	// TODO: Setup routes and middleware

	log.Printf("Server started on :%s\n", port)
	http.ListenAndServe(":"+port, nil)
}
