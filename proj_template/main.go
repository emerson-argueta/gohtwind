package main

import (
	"database/sql"
	"flag"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"{{PROJECT_NAME}}/infra"
)

func main() {
	envFileName := setUpEnv()
	err := godotenv.Load(envFileName)
	if err != nil {
		log.Fatal("Error loading env file")
	}
	port := os.Getenv("PORT")

	dbs, err := infra.NewDBsConfig().ConnectAll()
	if err != nil {
		log.Fatal(err)
	}
	for _, db := range dbs {
		defer func(db *sql.DB) {
			err := db.Close()
			if err != nil {
				log.Fatal(err)
			}
		}(db)
		err = db.Ping()
		if err != nil {
			log.Fatal(err)
		}
	}
	http.Handle("/static/", infra.LoggingMiddleware(http.StripPrefix("/static/", http.FileServer(http.Dir("./frontend/static/")))))
	// TODO: Setup routes and middleware

	log.Printf("Server started on :%s\n", port)
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func setUpEnv() string {
	// Define the -env flag
	env := flag.String("env", "development", "Environment (production or development)")
	flag.Parse()
	var envFileName string
	if *env == "production" {
		envFileName = ".env.production"
	} else {
		envFileName = ".env"
	}
	return envFileName
}
