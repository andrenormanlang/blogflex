package handlers

import (
    "encoding/json"
    "io"
    "log"
    "net/http"
    "net/url"
    "strconv"
    "strings"

    "github.com/a-h/templ"
    "github.com/gorilla/mux"

    "blogflex/internal/database"
    "blogflex/internal/models"
    "blogflex/views"
    "blogflex/internal/auth"
)

func CreateBlogHandler(w http.ResponseWriter, r *http.Request) {
    session := r.Context().Value("session").(*auth.Session)
    userID := session.UserID

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
    if err := database.DB.Where("blog_id = ?", blogID).Find(&posts).Error; err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    session, ok := r.Context().Value("session").(*auth.Session)
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }
    userID := session.UserID
    isOwner := userID == blog.UserID

    component := views.BlogPage(blog, posts, isOwner, true, "")
    templ.Handler(component).ServeHTTP(w, r)
}



