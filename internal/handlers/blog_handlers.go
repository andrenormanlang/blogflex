package handlers

import (
    "encoding/json"
    "io"
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
    "log"

)


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

    result, err := helpers.GraphQLRequest(query, variables)
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

    result, err := helpers.GraphQLRequest(query, nil)
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
                    id
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

    result, err := helpers.GraphQLRequest(query, variables)
    if err != nil {
        log.Printf("Failed to fetch blog: %v", err)
        http.Error(w, "Failed to fetch blog", http.StatusInternalServerError)
        return
    }

    log.Printf("GraphQL result: %+v\n", result)

    if result["data"] == nil {
        log.Println("result[data] is nil")
        http.Error(w, "Failed to fetch blog data", http.StatusInternalServerError)
        return
    }

    data, ok := result["data"].(map[string]interface{})
    if !ok || data["blogs_by_pk"] == nil {
        log.Println("data[blogs_by_pk] is nil or not a map")
        http.Error(w, "Failed to fetch blog data", http.StatusInternalServerError)
        return
    }

    blogData, ok := data["blogs_by_pk"].(map[string]interface{})
    if !ok {
        log.Println("blogData is nil or not a map")
        http.Error(w, "Failed to fetch blog data", http.StatusInternalServerError)
        return
    }

    var userMap map[string]interface{}
    if blogData["user"] != nil {
        userMap, ok = blogData["user"].(map[string]interface{})
        if !ok {
            userMap = map[string]interface{}{"id": "0", "username": "Unknown"}
        }
    } else {
        userMap = map[string]interface{}{"id": "0", "username": "Unknown"}
    }

    postsData, ok := blogData["posts"].([]interface{})
    if !ok {
        log.Println("postsData is nil or not a slice of interfaces")
        postsData = []interface{}{}
    }

    var posts []models.Post
    for _, postData := range postsData {
        postMap, ok := postData.(map[string]interface{})
        if !ok {
            log.Printf("postData is not a map: %v\n", postData)
            continue
        }

        var postUserMap map[string]interface{}
        if postMap["user"] != nil {
            postUserMap, ok = postMap["user"].(map[string]interface{})
            if !ok {
                postUserMap = map[string]interface{}{"username": "Unknown"}
            }
        } else {
            postUserMap = map[string]interface{}{"username": "Unknown"}
        }

        postID, ok := postMap["id"].(float64)
        if !ok {
            log.Printf("postMap[id] is not a float64: %v\n", postMap["id"])
            continue
        }

        posts = append(posts, models.Post{
            ID:       uint(postID),
            Title:    postMap["title"].(string),
            Content:  postMap["content"].(string),
            User: &models.User{
                Username: postUserMap["username"].(string),
            },
            FormattedCreatedAt: postMap["created_at"].(string),
        })
    }

    blogIDFloat, ok := blogData["id"].(float64)
    if !ok {
        log.Printf("blogData[id] is not a float64: %v\n", blogData["id"])
        http.Error(w, "Invalid blog ID", http.StatusInternalServerError)
        return
    }

    blog := models.Blog{
        ID:          uint(blogIDFloat),
        Name:        blogData["name"].(string),
        Description: blogData["description"].(string),
        User: &models.User{
            ID:       userMap["id"].(string),
            Username: userMap["username"].(string),
        },
        UserID: uint(blogIDFloat), // Convert to uint
    }

    // Check if user is authenticated to determine if they can create posts
    session, ok := r.Context().Value("session").(*sessions.Session)
    isOwner := false
    if ok {
        tokenStr, ok := session.Values["token"].(string)
        if ok {
            claims := &auth.Claims{}
            token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
                return auth.JwtKey, nil
            })
            if err == nil && token.Valid {
                isOwner = claims.UserID == blog.User.ID
            }
        }
    }

    loggedIn := helpers.IsLoggedIn(r)

    log.Println("Rendering blog page")
    component := views.BlogPage(blog, posts, isOwner, loggedIn, "")
    templ.Handler(component).ServeHTTP(w, r)
    log.Println("Blog page rendered successfully")
}
