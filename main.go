package main

import (
    "log"
    "net/http"
    "blogflex/internal/database"
    "blogflex/internal/router"
    "blogflex/internal/models"
)

func main() {
    // Initialize the database
    db := database.InitDatabase()
    
    // Automatically migrate the schema
    err := db.AutoMigrate(&models.User{}, &models.Post{}, &models.Comment{})
    if err != nil {
        log.Fatalf("Failed to migrate database schema: %v", err)
    }
 
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



// package main

// import (
//     "log"
//     "net/http"
//     "blogflex/internal/database"
//     "blogflex/internal/router"
//     "github.com/gorilla/sessions"
// )

// var store = sessions.NewCookieStore([]byte("your-very-secret-key"))

// func main() {
//     // Initialize the database
//     database.InitDatabase()

//     // Set up the router
//     r := router.SetupRouter(store)

//     // Serve static files
//     r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

//     // Start the server
//     log.Println("Server started at :8080")
//     if err := http.ListenAndServe(":8080", r); err != nil {
//         log.Fatal("ListenAndServe:", err)
//     }
// }

