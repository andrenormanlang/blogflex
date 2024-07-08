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
    "io"
    "net/url"
)

// CreatePostFormHandler handles the form submission for creating a post
func CreatePostFormHandler(w http.ResponseWriter, r *http.Request) {
    component := views.CreatePost() // Correctly refer to the templates.CreatePost component
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
// CreatePostHandler handles creating a new post
func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
    var post models.Post

    // Log the request body for debugging
    body, err := io.ReadAll(r.Body)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    log.Printf("Request Body: %s", body)

    // Determine content type
    contentType := r.Header.Get("Content-Type")

    if strings.Contains(contentType, "application/json") {
        // Decode JSON request body
        err = json.Unmarshal(body, &post)
        if err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }
    } else if strings.Contains(contentType, "application/x-www-form-urlencoded") {
        // Parse form-urlencoded request body
        values, err := url.ParseQuery(string(body))
        if err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        post.Title = values.Get("title")
        post.Content = values.Get("content")
        post.UserID = 1 // Hardcoded user ID for demonstration
    } else {
        http.Error(w, "Unsupported content type", http.StatusUnsupportedMediaType)
        return
    }

    result := database.DB.Create(&post)
    if result.Error != nil {
        http.Error(w, result.Error.Error(), http.StatusInternalServerError)
        return
    }

    // Respond with a success message
    w.Header().Set("Content-Type", "text/html")
    w.WriteHeader(http.StatusCreated)
    response := `<div class="bg-green-100 border-t border-b border-green-500 text-green-700 px-4 py-3" role="alert">
                    <p class="font-bold">Success!</p>
                    <p class="text-sm">Post created successfully.</p>
                 </div>`
    w.Write([]byte(response))
}