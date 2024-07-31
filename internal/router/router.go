package router

import (
    "blogflex/internal/handlers"
    "blogflex/internal/middleware"
    "net/http"

    "github.com/gorilla/mux"
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
    r.HandleFunc("/posts/{id}", handlers.PostDetailHandler).Methods("GET")
    r.HandleFunc("/posts/{id}/like", handlers.GetLikesCountHandler).Methods("GET")

    // Serve static files
    r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
    r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public/"))))

    // Protected routes
    protected := r.PathPrefix("/protected").Subrouter()
    protected.Use(middleware.AuthMiddleware)
    protected.HandleFunc("/logout", handlers.LogoutHandler).Methods("POST")
    protected.HandleFunc("/blogs/create", handlers.CreateBlogHandler).Methods("GET", "POST")
    protected.HandleFunc("/blogs/{id}/edit", handlers.EditBlogFormHandler).Methods("GET")
    protected.HandleFunc("/blogs/{id}/edit", handlers.EditBlogHandler).Methods("POST")
    protected.HandleFunc("/blogs/{id}", handlers.DeleteBlogHandler).Methods("DELETE")
    protected.HandleFunc("/posts/create", handlers.CreatePostFormHandler).Methods("GET")
    protected.HandleFunc("/posts", handlers.CreatePostHandler).Methods("POST")
    protected.HandleFunc("/posts/{id}/edit", handlers.EditPostFormHandler).Methods("GET")
    protected.HandleFunc("/posts/{id}/edit", handlers.EditPostHandler).Methods("POST")
    protected.HandleFunc("/posts/{id}", handlers.DeletePostHandler).Methods("DELETE")
    protected.HandleFunc("/posts/{id}/like", handlers.ToggleLikePostHandler).Methods("POST") // Ensure this route is defined
    protected.HandleFunc("/posts/{id}/comments", handlers.AddCommentHandler).Methods("POST")

    return r
}
