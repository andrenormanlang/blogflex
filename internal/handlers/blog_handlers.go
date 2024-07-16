package handlers

import (
    "encoding/json"
    "io"
    "bytes"
    "blogflex/internal/models"
    "net/http"
    "net/url"
    "strconv"
    "strings"
    "github.com/a-h/templ"
    "github.com/gorilla/mux"
    "blogflex/internal/auth"
    "blogflex/internal/helpers"
    "blogflex/views"
    "github.com/gorilla/sessions"
    "github.com/dgrijalva/jwt-go"
    "fmt"
    "log"
    "blogflex/internal/database"
)

func graphqlRequest(query string, variables map[string]interface{}) (map[string]interface{}, error) {
    requestBody, err := json.Marshal(database.GraphQLRequest{
        Query:     query,
        Variables: variables,
    })
    if err != nil {
        return nil, fmt.Errorf("failed to marshal GraphQL request: %v", err)
    }

    req, err := http.NewRequest("POST", database.HasuraEndpoint, bytes.NewBuffer(requestBody))
    if err != nil {
        return nil, fmt.Errorf("failed to create new HTTP request: %v", err)
    }
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("x-hasura-admin-secret", database.HasuraAdminSecret)

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("failed to perform HTTP request: %v", err)
    }
    defer resp.Body.Close()

    var result map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, fmt.Errorf("failed to decode GraphQL response: %v", err)
    }

    if errors, ok := result["errors"].([]interface{}); ok {
        var errorMessages []string
        for _, err := range errors {
            errorMap := err.(map[string]interface{})
            errorMessages = append(errorMessages, errorMap["message"].(string))
        }
        return nil, fmt.Errorf("GraphQL errors: %s", strings.Join(errorMessages, "; "))
    }

    return result, nil
}


// CreateBlogHandler handles creating a new blog
func CreateBlogHandler(w http.ResponseWriter, r *http.Request) {
    session := r.Context().Value("session").(*sessions.Session)
    userID := session.Values["userID"].(string)

    var blog struct {
        Name        string `json:"name"`
        Description string `json:"description"`
    }

    body, err := io.ReadAll(r.Body)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

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

    query := `
        mutation CreateBlog($name: String!, $description: String!, $user_id: uuid!) {
            insert_blogs_one(object: {name: $name, description: $description, user_id: $user_id}) {
                id
            }
        }
    `
    variables := map[string]interface{}{
        "name":        blog.Name,
        "description": blog.Description,
        "user_id":     userID,
    }

    result, err := graphqlRequest(query, variables)
    if err != nil || len(result["errors"].([]interface{})) > 0 {
        http.Error(w, "Failed to create blog", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    w.Header().Set("HX-Redirect", "/")
}

// BlogListHandler handles displaying a list of blogs
func BlogListHandler(w http.ResponseWriter, r *http.Request) {
    query := `
        query {
            blogs {
                id
                name
                description
                user {
                    username
                }
            }
        }
    `

    result, err := graphqlRequest(query, nil)
    if err != nil {
        log.Printf("Error fetching blogs: %v", err)
        http.Error(w, "Failed to fetch blogs: "+err.Error(), http.StatusInternalServerError)
        return
    }

    blogsData := result["data"].(map[string]interface{})["blogs"].([]interface{})
    var blogs []models.Blog
    for _, blogData := range blogsData {
        blogMap := blogData.(map[string]interface{})
        userMap := blogMap["user"].(map[string]interface{})

        blogs = append(blogs, models.Blog{
            ID:          uint(blogMap["id"].(float64)),
            Name:        blogMap["name"].(string),
            Description: blogMap["description"].(string),
            User: &models.User{
                Username: userMap["username"].(string),
            },
        })
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

    query := `
        query GetBlog($id: Int!) {
            blogs_by_pk(id: $id) {
                id
                name
                description
                user {
                    username
                }
                posts {
                    id
                    title
                    content
                    user {
                        username
                    }
                    created_at
                }
            }
        }
    `
    variables := map[string]interface{}{
        "id": blogID,
    }

    result, err := graphqlRequest(query, variables)
    if err != nil || len(result["errors"].([]interface{})) > 0 {
        http.Error(w, "Failed to fetch blog", http.StatusInternalServerError)
        return
    }

    // Ensure blogData and other fields exist
    blogData, ok := result["data"].(map[string]interface{})["blogs_by_pk"].(map[string]interface{})
    if !ok {
        http.Error(w, "Failed to parse blog data", http.StatusInternalServerError)
        return
    }

    var userMap map[string]interface{}
    if blogData["user"] != nil {
        userMap = blogData["user"].(map[string]interface{})
    } else {
        userMap = map[string]interface{}{"username": "Unknown"}
    }

    postsData, ok := blogData["posts"].([]interface{})
    if !ok {
        postsData = []interface{}{}
    }

    var posts []models.Post
    for _, postData := range postsData {
        postMap := postData.(map[string]interface{})
        var postUserMap map[string]interface{}
        if postMap["user"] != nil {
            postUserMap = postMap["user"].(map[string]interface{})
        } else {
            postUserMap = map[string]interface{}{"username": "Unknown"}
        }

        posts = append(posts, models.Post{
            ID:       uint(postMap["id"].(float64)),
            Title:    postMap["title"].(string),
            Content:  postMap["content"].(string),
            UserID:   uint(postMap["user_id"].(float64)),
            User: &models.User{
                Username: postUserMap["username"].(string),
            },
            FormattedCreatedAt: postMap["created_at"].(string),
        })
    }

    blog := models.Blog{
        ID:          uint(blogData["id"].(float64)),
        Name:        blogData["name"].(string),
        Description: blogData["description"].(string),
        User: &models.User{
            Username: userMap["username"].(string),
        },
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