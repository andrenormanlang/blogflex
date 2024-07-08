package handlers

import (
    "net/http"
    "github.com/a-h/templ"
    "blogflex/internal/database"
    "blogflex/internal/models"
    "encoding/json"
    "github.com/gorilla/mux"
    "strconv"
    "strings"
    "log"
    "blogflex/views" 
)

// CreatePostFormHandler handles the form submission for creating a post
func CreatePostFormHandler(w http.ResponseWriter, r *http.Request) {
    component := views.CreatePost() // Correctly refer to the templates.CreatePost component
    templ.Handler(component).ServeHTTP(w, r)
}

// PostListHandler handles displaying a list of posts
// PostListHandler handles displaying a list of posts
func PostListHandler(w http.ResponseWriter, r *http.Request) {
    var posts []models.Post
    result := database.DB.Find(&posts)
    if result.Error != nil {
        http.Error(w, result.Error.Error(), http.StatusInternalServerError)
        return
    }

    // Log the posts to ensure they have IDs
    for _, post := range posts {
        log.Printf("Post ID: %d, Title: %s", post.ID, post.Title)
    }

    // Check the Accept header to determine the response format
    acceptHeader := r.Header.Get("Accept")
    log.Printf("Accept header: %s", acceptHeader) // Log the Accept header value

    if strings.Contains(acceptHeader, "application/json") {
        log.Println("Returning JSON response") // Log the branch being executed
        // Set the content type to JSON
        w.Header().Set("Content-Type", "application/json")
        // Encode the posts as JSON and write to the response
        if err := json.NewEncoder(w).Encode(posts); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
    } else {
        log.Println("Returning HTML response") // Log the branch being executed
        // Render HTML template
        component := views.PostList(posts) // Correctly refer to the templates.PostList component
        templ.Handler(component).ServeHTTP(w, r)
    }
}
// PostDetailHandler handles displaying the details of a single post
func PostDetailHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid post ID", http.StatusBadRequest)
        log.Printf("Invalid post ID: %v", vars["id"]) // Log the invalid ID
        return
    }

    var post models.Post
    result := database.DB.First(&post, id)
    if result.Error != nil {
        http.Error(w, result.Error.Error(), http.StatusNotFound)
        log.Printf("Post not found: %v", id) // Log the ID not found
        return
    }

    log.Printf("Post found: ID=%d, Title=%s", post.ID, post.Title) // Log the found post

    // Check the Accept header to determine the response format
    acceptHeader := r.Header.Get("Accept")
    log.Printf("Accept header: %s", acceptHeader) // Log the Accept header value

    if strings.Contains(acceptHeader, "application/json") {
        log.Println("Returning JSON response") // Log the branch being executed
        // Set the content type to JSON
        w.Header().Set("Content-Type", "application/json")
        // Encode the post as JSON and write to the response
        if err := json.NewEncoder(w).Encode(post); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
    } else {
        log.Println("Returning HTML response") // Log the branch being executed
        // Render HTML template
        component := views.PostDetail(post) // Correctly refer to the templates.PostDetail component
        templ.Handler(component).ServeHTTP(w, r)
    }
}
// CreatePostHandler handles creating a new post
func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
    var post models.Post
    err := json.NewDecoder(r.Body).Decode(&post)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    result := database.DB.Create(&post)
    if result.Error != nil {
        http.Error(w, result.Error.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(post)
}
