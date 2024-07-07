package main

import (
    "log"
    "net/http"
    "blogflex/internal/database"
    "blogflex/internal/router"
)

func main() {
    // Initialize the database
    database.InitDatabase()

    // Set up the router
    r := router.SetupRouter()

    // Serve static files
    r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

    // Start the server
    log.Println("Server started at :8080")
    if err := http.ListenAndServe(":8080", r); err != nil {
        log.Fatal("ListenAndServe:", err)
    }
}
