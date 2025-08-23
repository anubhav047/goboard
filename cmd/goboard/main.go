package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	"github.com/anubhav047/goboard/internal/db"
	httphandlers "github.com/anubhav047/goboard/internal/http"
	boardservice "github.com/anubhav047/goboard/internal/services/board"
	cardservice "github.com/anubhav047/goboard/internal/services/card"
	listservice "github.com/anubhav047/goboard/internal/services/list"
	userservice "github.com/anubhav047/goboard/internal/services/user"
	_ "github.com/jackc/pgx/v5/stdlib"
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

	// CREATE STANDARD DB POOL FOR SCS
	// We need this because postgresstore expects a *sql.DB object.
	dbConn, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatalf("Unable to create standard db connection pool: %v\n", err)
	}
	defer dbConn.Close()
	// Ping to confirm the connection is alive.
	if err := dbConn.Ping(); err != nil {
		log.Fatalf("Unable to ping database: %v\n", err)
	}

	//Connect to the database
	dbpool, err := db.Connect(connStr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	//Defer the closing of the connection pool
	defer dbpool.Close()

	log.Println("Database connection successful.")

	// SESSION MANAGER
	sessionManager := scs.New()
	sessionManager.Store = postgresstore.New(dbConn)
	sessionManager.Lifetime = 24 * time.Hour
	sessionManager.Cookie.Persist = true
	sessionManager.Cookie.Secure = false
	sessionManager.Cookie.SameSite = http.SameSiteLaxMode

	// Create a Queries object from the connection pool
	queries := db.New(dbpool)

	// Create the user Service
	userService := userservice.New(queries)

	// Create the board Service
	boardService := boardservice.New(queries)

	// Create the list Service
	listService := listservice.New(queries)

	// Create the card Service
	cardService := cardservice.New(queries)

	// Create middleware struct
	mw := httphandlers.NewMiddleware(sessionManager, queries)

	// Create and register User Handler
	userHandler := httphandlers.NewUserHandler(userService, sessionManager)

	// Create and register Board Handler
	boardHandler := httphandlers.NewBoardHandler(boardService)

	// Create and register List Handler
	listHandler := httphandlers.NewListHandler(listService)

	// Create and register Card Handler
	cardHandler := httphandlers.NewCardHandler(cardService)

	mux := http.NewServeMux()
	userHandler.RegisterRoutes(mux, mw)
	boardHandler.RegisterRoutes(mux, mw)
	listHandler.RegisterRoutes(mux, mw)
	cardHandler.RegisterRoutes(mux, mw)
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok")
	})

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", sessionManager.LoadAndSave(mux)))

}
