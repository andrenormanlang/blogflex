package handlers

import (
    "encoding/json"
    "io"
    "log"
    "net/http"
    "net/url"
    "strings"
    "github.com/gorilla/sessions"
    "blogflex/internal/database"
    "blogflex/internal/models"
)

var store = sessions.NewCookieStore([]byte("super-secret-key"))

func LoginHandler(w http.ResponseWriter, r *http.Request) {
    var user models.User

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
        err = json.Unmarshal(body, &user)
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

        user.Username = values.Get("username")
        user.Password = values.Get("password")
    } else {
        http.Error(w, "Unsupported content type", http.StatusUnsupportedMediaType)
        return
    }

    var dbUser models.User
    result := database.DB.Where("username = ? AND password = ?", user.Username, user.Password).First(&dbUser)
    if result.Error != nil {
        http.Error(w, "Invalid username or password", http.StatusUnauthorized)
        return
    }

    session, _ := store.Get(r, "session-name")
    session.Values["userID"] = dbUser.ID
    session.Save(r, w)

    log.Printf("User logged in: %s, session: %v", dbUser.Username, session.Values)

    // Redirect to the blog page
    http.Redirect(w, r, "/protected/posts", http.StatusFound)
}
