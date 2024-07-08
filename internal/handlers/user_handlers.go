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
)


func ListUsersHandler(w http.ResponseWriter, r *http.Request) {
    var users []models.User
    result := database.DB.Find(&users)
    if result.Error != nil {
        http.Error(w, result.Error.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(users)
}

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }

    var user models.User
    result := database.DB.First(&user, id)
    if result.Error != nil {
        http.Error(w, result.Error.Error(), http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
    var user models.User
    err := json.NewDecoder(r.Body).Decode(&user)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    result := database.DB.Create(&user)
    if result.Error != nil {
        http.Error(w, result.Error.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}

func MainPageHandler(w http.ResponseWriter, r *http.Request) {
    component := views.MainPage()
    templ.Handler(component).ServeHTTP(w, r)
}

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
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
        user.Email = values.Get("email")
        user.Password = values.Get("password")
    } else {
        http.Error(w, "Unsupported content type", http.StatusUnsupportedMediaType)
        return
    }

    result := database.DB.Create(&user)
    if result.Error != nil {
        http.Error(w, result.Error.Error(), http.StatusInternalServerError)
        return
    }

    // Set session values
    session, _ := store.Get(r, "session-name")
    session.Values["userID"] = user.ID
    session.Save(r, w)

    // Redirect to the blog page
    w.Header().Set("HX-Redirect", "/protected/posts")
    w.WriteHeader(http.StatusCreated)
    response := map[string]string{"message": "Sign-up successful"}
    json.NewEncoder(w).Encode(response)
}

