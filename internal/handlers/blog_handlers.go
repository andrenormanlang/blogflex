package handlers

import (
    "encoding/json"
    "io"
    "net/http"
    "log"
    "blogflex/internal/database"
    "blogflex/internal/models"
    "github.com/gorilla/sessions"
    "blogflex/views"
	"github.com/a-h/templ"
	"strings"
	"net/url"
    "strconv"
    "github.com/gorilla/mux"
)

// CreateBlogHandler handles creating a new blog
func CreateBlogHandler(w http.ResponseWriter, r *http.Request) {
    session := r.Context().Value("session").(*sessions.Session)
    userID, ok := session.Values["userID"].(uint)
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    var blog models.Blog
    body, err := io.ReadAll(r.Body)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    log.Printf("Request Body: %s", body)

    contentType := r.Header.Get("Content-Type")
    if strings.Contains(contentType, "application/json") {
        err = json.Unmarshal(body, &blog)
        if err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }
    } else if strings.Contains(contentType, "application/x-www-form-urlencoded") {
        values, err := url.ParseQuery(string(body))
        if err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }
        blog.Name = values.Get("name")
        blog.Description = values.Get("description")
    } else {
        http.Error(w, "Unsupported content type", http.StatusUnsupportedMediaType)
        return
    }

    blog.UserID = userID // Set the user ID from session

    result := database.DB.Create(&blog)
    if result.Error != nil {
        http.Error(w, result.Error.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    w.Header().Set("HX-Redirect", "/")
}

// BlogListHandler handles displaying a list of blogs
func BlogListHandler(w http.ResponseWriter, r *http.Request) {
    var blogs []models.Blog
    result := database.DB.Preload("User").Find(&blogs)
    if result.Error != nil {
        http.Error(w, result.Error.Error(), http.StatusInternalServerError)
        return
    }

    component := views.BlogList(blogs)
    templ.Handler(component).ServeHTTP(w, r)
}

func BlogPageHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    blogID, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid blog ID", http.StatusBadRequest)
        return
    }

    var blog models.Blog
    if err := database.DB.Preload("User").First(&blog, blogID).Error; err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }

    var posts []models.Post
    if err := database.DB.Where("blog_id = ?", blogID).Preload("User").Find(&posts).Error; err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    component := views.BlogPage(blog, posts)
    templ.Handler(component).ServeHTTP(w, r)
}

