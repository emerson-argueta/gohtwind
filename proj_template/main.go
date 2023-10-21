package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"{{PROJECT_NAME}}/infra"
)

var db *sql.DB

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not specified
	}
	var (
		host     = os.Getenv("DB_HOST")
		dbport   = os.Getenv("DB_PORT")
		user     = os.Getenv("DB_USER")
		password = os.Getenv("DB_PASSWORD")
		dbName   = os.Getenv("DB_NAME")
	)
	connectString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", user, password, host, dbport, dbName)
	db, err = sql.Open("mysql", connectString)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	http.Handle("/static/", infra.LoggingMiddleware(http.StripPrefix("/static/", http.FileServer(http.Dir("./frontend/static/")))))
	// TODO: Setup routes and middleware

	log.Printf("Server started on :%s\n", port)
	http.ListenAndServe(":"+port, nil)
}
