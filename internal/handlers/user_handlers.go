package handlers

import (

    "encoding/json"
    "io"
    "log"
    "net/http"
    "net/url"
    "strconv"
    "strings"
    "time"

    "github.com/a-h/templ"
    "github.com/gorilla/mux"
    "github.com/dgrijalva/jwt-go"

    "blogflex/internal/auth"
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
    var blogs []models.Blog
    result := database.DB.Preload("User").Find(&blogs)
    if result.Error != nil {
        http.Error(w, result.Error.Error(), http.StatusInternalServerError)
        return
    }

    component := views.MainPage(blogs)
    templ.Handler(component).ServeHTTP(w, r)
}

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
    var user models.User
    var blog models.Blog

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
        blog.Name = values.Get("blog-name")
        blog.Description = values.Get("blog-description")
    } else {
        http.Error(w, "Unsupported content type", http.StatusUnsupportedMediaType)
        return
    }

    result := database.DB.Create(&user)
    if result.Error != nil {
        http.Error(w, result.Error.Error(), http.StatusInternalServerError)
        return
    }

    blog.UserID = user.ID

    result = database.DB.Create(&blog)
    if result.Error != nil {
        http.Error(w, result.Error.Error(), http.StatusInternalServerError)
        return
    }

    expirationTime := time.Now().Add(5 * time.Minute)
    claims := &auth.Claims{
        UserID:   user.ID,
        Username: user.Username,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: expirationTime.Unix(),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString(auth.JwtKey)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    http.SetCookie(w, &http.Cookie{
        Name:    "token",
        Value:   tokenString,
        Expires: expirationTime,
    })

    log.Printf("User signed up: %s", user.Username)
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "redirect": "/blogs/" + strconv.Itoa(int(blog.ID)),
    })
}


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
    result := database.DB.Where("username = ?", user.Username).First(&dbUser)
    if result.Error != nil || dbUser.Password != user.Password {
        http.Error(w, "Invalid username or password", http.StatusUnauthorized)
        return
    }

    expirationTime := time.Now().Add(5 * time.Minute)
    claims := &auth.Claims{
        UserID:   dbUser.ID,
        Username: dbUser.Username,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: expirationTime.Unix(),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString(auth.JwtKey)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    http.SetCookie(w, &http.Cookie{
        Name:    "token",
        Value:   tokenString,
        Expires: expirationTime,
    })

    // Fetch the user's blog
    var blog models.Blog
    result = database.DB.Where("user_id = ?", dbUser.ID).First(&blog)
    if result.Error != nil {
        http.Error(w, "User does not have a blog", http.StatusInternalServerError)
        return
    }

    log.Printf("User logged in: %s", dbUser.Username)
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "redirect": "/blogs/" + strconv.Itoa(int(blog.ID)),
    })
}




func LogoutHandler(w http.ResponseWriter, r *http.Request) {
    // Invalidate the session token
    cookie := &http.Cookie{
        Name:     "token",
        Value:    "",
        Path:     "/",
        MaxAge:   -1,
        HttpOnly: true,
    }
    http.SetCookie(w, cookie)
    w.Header().Set("HX-Redirect", "/")
    w.WriteHeader(http.StatusOK)
}



// package handlers

// import (
//     "encoding/json"
//     "io"
//     "log"
//     "net/http"
//     "net/url"
//     "strconv"
//     "strings"

//     "github.com/a-h/templ"
//     "github.com/gorilla/mux"


//     "blogflex/internal/database"
//     "blogflex/internal/models"
//     "blogflex/views"
// 	"github.com/gorilla/sessions"
// )

// var store = sessions.NewCookieStore([]byte("your-very-secret-key"))

// func ListUsersHandler(w http.ResponseWriter, r *http.Request) {
//     var users []models.User
//     result := database.DB.Find(&users)
//     if result.Error != nil {
//         http.Error(w, result.Error.Error(), http.StatusInternalServerError)
//         return
//     }

//     w.Header().Set("Content-Type", "application/json")
//     json.NewEncoder(w).Encode(users)
// }

// func GetUserHandler(w http.ResponseWriter, r *http.Request) {
//     vars := mux.Vars(r)
//     id, err := strconv.Atoi(vars["id"])
//     if err != nil {
//         http.Error(w, "Invalid user ID", http.StatusBadRequest)
//         return
//     }

//     var user models.User
//     result := database.DB.First(&user, id)
//     if result.Error != nil {
//         http.Error(w, result.Error.Error(), http.StatusNotFound)
//         return
//     }

//     w.Header().Set("Content-Type", "application/json")
//     json.NewEncoder(w).Encode(user)
// }

// func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
//     var user models.User
//     err := json.NewDecoder(r.Body).Decode(&user)
//     if err != nil {
//         http.Error(w, err.Error(), http.StatusBadRequest)
//         return
//     }

//     result := database.DB.Create(&user)
//     if result.Error != nil {
//         http.Error(w, result.Error.Error(), http.StatusInternalServerError)
//         return
//     }

//     w.Header().Set("Content-Type", "application/json")
//     json.NewEncoder(w).Encode(user)
// }

// func MainPageHandler(w http.ResponseWriter, r *http.Request) {
//     component := views.MainPage()
//     templ.Handler(component).ServeHTTP(w, r)
// }

// func SignUpHandler(w http.ResponseWriter, r *http.Request) {
//     var user models.User

//     // Log the request body for debugging
//     body, err := io.ReadAll(r.Body)
//     if err != nil {
//         http.Error(w, err.Error(), http.StatusBadRequest)
//         return
//     }
//     log.Printf("Request Body: %s", body)

//     // Determine content type
//     contentType := r.Header.Get("Content-Type")

//     if strings.Contains(contentType, "application/json") {
//         // Decode JSON request body
//         err = json.Unmarshal(body, &user)
//         if err != nil {
//             http.Error(w, err.Error(), http.StatusBadRequest)
//             return
//         }
//     } else if strings.Contains(contentType, "application/x-www-form-urlencoded") {
//         // Parse form-urlencoded request body
//         values, err := url.ParseQuery(string(body))
//         if err != nil {
//             http.Error(w, err.Error(), http.StatusBadRequest)
//             return
//         }

//         user.Username = values.Get("username")
//         user.Email = values.Get("email")
//         user.Password = values.Get("password")
//     } else {
//         http.Error(w, "Unsupported content type", http.StatusUnsupportedMediaType)
//         return
//     }

//     result := database.DB.Create(&user)
//     if result.Error != nil {
//         http.Error(w, result.Error.Error(), http.StatusInternalServerError)
//         return
//     }

//     // Set session values
//     session, _ := store.Get(r, "session-name")
//     session.Values["userID"] = user.ID
//     session.Save(r, w)

//     // Redirect to the blog page
//     w.Header().Set("HX-Redirect", "/protected/posts")
//     w.WriteHeader(http.StatusCreated)
//     response := map[string]string{"message": "Sign-up successful"}
//     json.NewEncoder(w).Encode(response)
// }

