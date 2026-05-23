package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"Shazam-Vscode/backend/internal/api"
	"Shazam-Vscode/backend/internal/db"
)

func main() {
	_ = godotenv.Load()

	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("DB_DSN environment variable is required")
	}

	if err := db.Connect(dsn); err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer db.Close()

	router := api.NewRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}