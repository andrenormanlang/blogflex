package main

import (
    "log"
    "net/http"
    "blogflex/internal/database"
    "blogflex/internal/router"
    "github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("your-very-secret-key"))

func main() {
    // Initialize the database
    database.InitDatabase()

    // Set up the router
    r := router.SetupRouter(store)

    // Serve static files
    r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

    // Start the server
    log.Println("Server started at :8080")
    if err := http.ListenAndServe(":8080", r); err != nil {
        log.Fatal("ListenAndServe:", err)
    }
}

