package main

import (
    "log"
    "net/http"
    "os"
    "blogflex/internal/database"
    "blogflex/internal/router"
    "github.com/joho/godotenv"
)

func main() {
    // Load environment variables from .env file if it exists
    err := godotenv.Load()
    if err != nil {
        log.Println("No .env file found")
    }

    // Check required environment variables
    hasuraAdminSecret := os.Getenv("HASURA_ADMIN_SECRET")
    hasuraEndpoint := os.Getenv("HASURA_ENDPOINT")
    if hasuraAdminSecret == "" || hasuraEndpoint == "" {
        log.Fatal("Missing required environment variables: HASURA_ADMIN_SECRET or HASURA_ENDPOINT")
    }

    // Initialize the database with environment variables
    database.InitHasura()

    // Set up the router
    r := router.SetupRouter()

    // Start the server
    log.Println("Server started at :8080")
    if err := http.ListenAndServe(":8080", r); err != nil {
        log.Fatal("ListenAndServe:", err)
    }
}