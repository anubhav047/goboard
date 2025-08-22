package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/anubhav047/goboard/internal/db"
	"github.com/joho/godotenv"
)

func main() {
	//Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using environment variables")
	}

	//Get database connection string from environment variables
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	//Connect to the database
	dbpool, err := db.Connect(connStr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	//Defer the closing of the connection pool
	defer dbpool.Close()

	log.Println("Database connection successful.")

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok") //response for healthz endpoint
	})

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))

}
