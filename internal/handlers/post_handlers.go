package handlers

import (
    "net/http"
    "github.com/a-h/templ"
    "blogflex/internal/database"
    "blogflex/internal/models"
    "encoding/json"
    "github.com/gorilla/mux"
    "strconv"
    "blogflex/views" // Make sure this import path is correct
)

// CreatePostFormHandler handles the form submission for creating a post
func CreatePostFormHandler(w http.ResponseWriter, r *http.Request) {
    component := views.CreatePost() // Correctly refer to the templates.CreatePost component
    templ.Handler(component).ServeHTTP(w, r)
}

// PostDetailHandler handles displaying the details of a single post
func PostDetailHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid post ID", http.StatusBadRequest)
        return
    }

    var post models.Post
    result := database.DB.First(&post, id)
    if result.Error != nil {
        http.Error(w, result.Error.Error(), http.StatusNotFound)
        return
    }

    component := views.PostDetail(post) // Correctly refer to the templates.PostDetail component
    templ.Handler(component).ServeHTTP(w, r)
}

// PostListHandler handles displaying a list of posts
func PostListHandler(w http.ResponseWriter, r *http.Request) {
    var posts []models.Post
    result := database.DB.Find(&posts)
    if result.Error != nil {
        http.Error(w, result.Error.Error(), http.StatusInternalServerError)
        return
    }

    component := views.PostList(posts) // Ensure PostList is correctly defined
    templ.Handler(component).ServeHTTP(w, r)
}


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