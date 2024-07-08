package router

import (
    "context"
    "net/http"

    "github.com/gorilla/mux"
    "github.com/gorilla/sessions"
    "blogflex/internal/handlers"
    "blogflex/middleware"
)

func SetupRouter(store *sessions.CookieStore) *mux.Router {
    r := mux.NewRouter()

    // Middleware to handle sessions
    r.Use(func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            session, _ := store.Get(r, "session-name")
            r = r.WithContext(context.WithValue(r.Context(), "session", session))
            next.ServeHTTP(w, r)
        })
    })

    // Public routes
    r.HandleFunc("/", handlers.MainPageHandler).Methods("GET")
    r.HandleFunc("/signup", handlers.SignUpHandler).Methods("POST")
    r.HandleFunc("/login", handlers.LoginHandler).Methods("POST")

    // Protected routes
    protected := r.PathPrefix("/protected").Subrouter()
    protected.Use(middleware.AuthMiddleware(store))
    protected.HandleFunc("/posts", handlers.PostListHandler).Methods("GET")
    protected.HandleFunc("/posts/create", handlers.CreatePostFormHandler).Methods("GET")
    protected.HandleFunc("/posts/create", handlers.CreatePostHandler).Methods("POST")
    protected.HandleFunc("/posts/{id}", handlers.PostDetailHandler).Methods("GET")
    protected.HandleFunc("/users/{id}", handlers.GetUserHandler).Methods("GET")

    // Serve static files
    r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

    return r
}

