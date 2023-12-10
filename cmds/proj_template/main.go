package main

import (
	"database/sql"
	"flag"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"{{PROJECT_NAME}}/auth"
	"{{PROJECT_NAME}}/infra"
)

func main() {
	setUpEnv()
	dbs := setUpDBs()
	setUpRoutes(dbs)
	setUpServer(dbs)
}

func setUpEnv() {
	env := flag.String("env", "development", "Environment (production or development)")
	flag.Parse()
	var envFileName string
	if *env == "production" {
		return
	}
	envFileName = ".env"
	err := godotenv.Load(envFileName)
	if err != nil {
		log.Fatal("Error loading env file")
	}
}

func setUpDBs() map[string]*sql.DB {
	dbs, err := infra.NewDBsConfig().ConnectAll()
	if err != nil {
		log.Fatal(err)
	}
	for _, db := range dbs {
		infra.CheckDatabaseConnection(db)
	}
	return dbs
}

func setUpRoutes(dbs map[string]*sql.DB) {
	http.Handle("/static/", infra.LoggingMiddleware(http.StripPrefix("/static/", http.FileServer(http.Dir("./frontend/static/")))))
	auth.SetupRoutes(dbs, infra.LoggingMiddleware)
	// TODO: Add more routes here

	// Then activate the routes
	infra.ActivateRoutes()
}

func setUpServer(dbs map[string]*sql.DB) {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		log.Printf("Received signal: %s", sig)
		done <- true
	}()

	go func() {
		port := os.Getenv("PORT")
		log.Printf("Server started on :%s\n", port)
		if err := http.ListenAndServe(":"+port, nil); err != nil {
			log.Fatal(err)
		}
	}()
	<-done
	log.Println("Shutting down server...")

	for name, db := range dbs {
		if err := db.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
		log.Printf("Database %s closed", name)
	}

	log.Println("Server gracefully stopped")
}
