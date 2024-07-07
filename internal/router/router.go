package router

import (
    "github.com/gorilla/mux"
    "blogflex/internal/handlers"
)

// SetupRouter initializes the router with the defined routes.
func SetupRouter() *mux.Router {
    r := mux.NewRouter()

    // Route to display the form for creating a new post
    r.HandleFunc("/posts/create", handlers.CreatePostFormHandler).Methods("GET")

    // Route to handle the creation of a new post
    r.HandleFunc("/posts", handlers.CreatePostHandler).Methods("POST")

    // Route to display the list of posts
    r.HandleFunc("/posts", handlers.PostListHandler).Methods("GET")

    // Route to display the details of a single post
    r.HandleFunc("/posts/{id}", handlers.PostDetailHandler).Methods("GET")

    // User-related routes
    r.HandleFunc("/users", handlers.CreateUserHandler).Methods("POST")
    r.HandleFunc("/users", handlers.ListUsersHandler).Methods("GET")
    r.HandleFunc("/users/{id}", handlers.GetUserHandler).Methods("GET")

    // Comment routes
    r.HandleFunc("/posts/{postID}/comments", handlers.CreateCommentHandler).Methods("POST")
    r.HandleFunc("/posts/{postID}/comments", handlers.GetCommentsHandler).Methods("GET")

    return r
}



