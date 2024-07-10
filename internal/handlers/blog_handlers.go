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
    "github.com/gorilla/sessions"
    "github.com/dgrijalva/jwt-go"
    "blogflex/internal/helpers"
    
)

// CreateBlogHandler handles creating a new blog
// CreateBlogHandler handles creating a new blog
func CreateBlogHandler(w http.ResponseWriter, r *http.Request) {
    session := r.Context().Value("session").(*sessions.Session)
    userID := session.Values["userID"].(uint)

    // Check if the user already has a blog
    var existingBlog models.Blog
    if err := database.DB.Where("user_id = ?", userID).First(&existingBlog).Error; err == nil {
        http.Error(w, "You already have a blog", http.StatusBadRequest)
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

// BlogPageHandler handles displaying a single blog page with its posts
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

    session, ok := r.Context().Value("session").(*sessions.Session)
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }
    tokenStr, ok := session.Values["token"].(string)
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    claims := &auth.Claims{}
    token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
        return auth.JwtKey, nil
    })
    if err != nil || !token.Valid {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    userID := claims.UserID
    isOwner := userID == blog.UserID
    loggedIn := helpers.IsLoggedIn(r)

    component := views.BlogPage(blog, posts, isOwner, loggedIn, "")
    templ.Handler(component).ServeHTTP(w, r)
}


