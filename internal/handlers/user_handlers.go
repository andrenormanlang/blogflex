package handlers

import (
    "encoding/json"
    "fmt"
    "io"
    "log"
    "net/http"
    "net/url"
    "strconv"
    "strings"
    "time"

    "github.com/dgrijalva/jwt-go"
    "github.com/gorilla/mux"
    "github.com/gorilla/sessions"
    "golang.org/x/crypto/bcrypt"
    "blogflex/internal/database"
    "github.com/a-h/templ"
    "blogflex/views"
   

    "blogflex/internal/auth"
    "blogflex/internal/models"
    "blogflex/internal/helpers"
)

var store = sessions.NewCookieStore([]byte("your-very-secret-key"))

func MainPageHandler(w http.ResponseWriter, r *http.Request) {
    query := `
        query {
            blogs {
                id
                name
                description
                user {
                    username
                }
                created_at
            }
        }
    `
    result, err := database.ExecuteGraphQL(query, nil)
    if err != nil {
        http.Error(w, "Failed to fetch blogs", http.StatusInternalServerError)
        return
    }

    blogsData := result["blogs"].([]interface{})
    var blogs []models.Blog
    for _, blogData := range blogsData {
        blogMap := blogData.(map[string]interface{})
        userMap := blogMap["user"].(map[string]interface{})
        blogs = append(blogs, models.Blog{
            ID:                 uint(blogMap["id"].(float64)),
            Name:               blogMap["name"].(string),
            Description:        blogMap["description"].(string),
            FormattedCreatedAt: blogMap["created_at"].(string),
            User: &models.User{
                Username: userMap["username"].(string),
            },
        })
    }

    component := views.MainPage(blogs, false) // Adjust based on your logic for loggedIn
    templ.Handler(component).ServeHTTP(w, r)
}

// ListUsersHandler handles listing all users
func ListUsersHandler(w http.ResponseWriter, r *http.Request) {
    query := `
        query {
            users {
                id
                username
                email
            }
        }
    `
    result, err := helpers.GraphQLRequest(query, nil)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    users := result["data"].(map[string]interface{})["users"].([]interface{})
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(users)
}

// GetUserHandler handles fetching a single user by ID
func GetUserHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }

    query := `
        query ($id: Int!) {
            users_by_pk(id: $id) {
                id
                username
                email
            }
        }
    `
    variables := map[string]interface{}{
        "id": id,
    }

    result, err := helpers.GraphQLRequest(query, variables)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    user := result["data"].(map[string]interface{})["users_by_pk"]
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}

// SignUpHandler handles user registration
// SignUpHandler handles user registration

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
    var user models.User
    var blogName, blogDescription string

    body, err := io.ReadAll(r.Body)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    log.Printf("Request Body: %s", body)

    contentType := r.Header.Get("Content-Type")

    if strings.Contains(contentType, "application/json") {
        err = json.Unmarshal(body, &user)
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

        user.Username = values.Get("username")
        user.Email = values.Get("email")
        user.Password = values.Get("password")
        blogName = values.Get("blogName")
        blogDescription = values.Get("blogDescription")
    } else {
        http.Error(w, "Unsupported content type", http.StatusUnsupportedMediaType)
        return
    }

    // Hash the password before storing it
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    user.Password = string(hashedPassword)

    mutation := `
        mutation CreateUser($username: String!, $email: String!, $password: String!) {
            insert_users_one(object: {username: $username, email: $email, password: $password}) {
                id
                username
                email
            }
        }
    `
    variables := map[string]interface{}{
        "username": user.Username,
        "email":    user.Email,
        "password": user.Password,
    }

    result, err := helpers.GraphQLRequest(mutation, variables)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Create a blog for the new user
    userId := result["data"].(map[string]interface{})["insert_users_one"].(map[string]interface{})["id"].(string)
    blogMutation := `
        mutation CreateBlog($user_id: uuid!, $name: String!, $description: String!) {
            insert_blogs_one(object: {user_id: $user_id, name: $name, description: $description}) {
                id
            }
        }
    `
    blogVariables := map[string]interface{}{
        "user_id":     userId,
        "name":        blogName,
        "description": blogDescription,
    }

    _, err = helpers.GraphQLRequest(blogMutation, blogVariables)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    log.Printf("User signed up: %s", user.Username)
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "message": "Sign up successful! Please log in to continue.",
    })
}
// LoginHandler handles user login
func LoginHandler(w http.ResponseWriter, r *http.Request) {
    var user models.User

    body, err := io.ReadAll(r.Body)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    log.Printf("Request Body: %s", body)

    contentType := r.Header.Get("Content-Type")

    if strings.Contains(contentType, "application/json") {
        err = json.Unmarshal(body, &user)
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

        user.Username = values.Get("username")
        user.Password = values.Get("password")
    } else {
        http.Error(w, "Unsupported content type", http.StatusUnsupportedMediaType)
        return
    }

    query := `
        query GetUser($username: String!) {
            users(where: {username: {_eq: $username}}) {
                id
                username
                password
            }
        }
    `
    variables := map[string]interface{}{
        "username": user.Username,
    }

    result, err := helpers.GraphQLRequest(query, variables)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    users := result["data"].(map[string]interface{})["users"].([]interface{})
    if len(users) == 0 {
        http.Error(w, "Invalid username or password", http.StatusUnauthorized)
        return
    }

    dbUser := users[0].(map[string]interface{})

    // Compare the stored hashed password with the provided password
    err = bcrypt.CompareHashAndPassword([]byte(dbUser["password"].(string)), []byte(user.Password))
    if err != nil {
        http.Error(w, "Invalid username or password", http.StatusUnauthorized)
        return
    }

    // Convert the id to a string
    userID := fmt.Sprintf("%v", dbUser["id"])

    expirationTime := time.Now().Add(24 * time.Hour)
    claims := &auth.Claims{
        UserID:   userID,
        Username: dbUser["username"].(string),
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

    session, _ := store.Get(r, "session-name")
    session.Values["token"] = tokenString
    session.Values["userID"] = userID // Store the user ID in the session
    err = session.Save(r, w)
    if err != nil {
        http.Error(w, "Failed to save session", http.StatusInternalServerError)
        return
    }

    log.Printf("User logged in: %s", dbUser["username"])

    blogQuery := `
        query GetUserBlog($user_id: uuid!) {
            blogs(where: {user_id: {_eq: $user_id}}) {
                id
            }
        }
    `
    blogVariables := map[string]interface{}{
        "user_id": userID,
    }

    blogResult, err := helpers.GraphQLRequest(blogQuery, blogVariables)
    if err != nil {
        http.Error(w, "User does not have a blog", http.StatusBadRequest)
        return
    }

    blogs := blogResult["data"].(map[string]interface{})["blogs"].([]interface{})
    if len(blogs) == 0 {
        http.Error(w, "User does not have a blog", http.StatusBadRequest)
        return
    }

    blogID := fmt.Sprintf("%v", blogs[0].(map[string]interface{})["id"])

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "redirect": fmt.Sprintf("/blogs/%s", blogID),
    })
}

// LogoutHandler handles user logout
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
    session, _ := store.Get(r, "session-name")
    session.Values["token"] = ""
    session.Save(r, w)

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


