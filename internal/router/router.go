package router

import (
    "github.com/gorilla/mux"
    "blogflex/internal/handlers"
    "blogflex/internal/middleware"
    "net/http"
)

func SetupRouter() *mux.Router {
    r := mux.NewRouter()

    // Apply session middleware
    r.Use(middleware.SessionMiddleware)

    // Public routes
    r.HandleFunc("/", handlers.MainPageHandler).Methods("GET")
    r.HandleFunc("/signup", handlers.SignUpHandler).Methods("POST")
    r.HandleFunc("/login", handlers.LoginHandler).Methods("POST")
    r.HandleFunc("/blogs/{id}", handlers.BlogPageHandler).Methods("GET")
    r.HandleFunc("/posts/{id}", handlers.PostDetailHandler).Methods("GET") // Moved outside protected routes
    
    // Serve static files
    r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

    // Protected routes
    protected := r.PathPrefix("/protected").Subrouter()
    protected.Use(middleware.AuthMiddleware)
    protected.HandleFunc("/logout", handlers.LogoutHandler).Methods("POST")
    protected.HandleFunc("/blogs/create", handlers.CreateBlogHandler).Methods("GET", "POST")
    protected.HandleFunc("/posts/create", handlers.CreatePostFormHandler).Methods("GET")
    protected.HandleFunc("/posts", handlers.CreatePostHandler).Methods("POST")  // Ensure this is correct
    // protected.HandleFunc("/posts/{id}", handlers.PostDetailHandler).Methods("GET")

    return r
}

