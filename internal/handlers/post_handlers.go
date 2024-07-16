package handlers

import (
    "net/http"
    "github.com/a-h/templ"
    "blogflex/views"
    // "github.com/gorilla/sessions"
    "encoding/json"
    "io"
    "net/url"
    "strings"
    "log"
    "blogflex/internal/models"
    "strconv"
    "github.com/gorilla/mux"
    "blogflex/internal/database"
    "fmt"
)

func CreatePostFormHandler(w http.ResponseWriter, r *http.Request) {
    session, _ := store.Get(r, "session-name")
    userID, ok := session.Values["userID"].(string)
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    component := views.CreatePost(userID)
    templ.Handler(component).ServeHTTP(w, r)
}



func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
    session, _ := store.Get(r, "session-name")
    userID, ok := session.Values["userID"].(string)
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    var post struct {
        Title   string `json:"title"`
        Content string `json:"content"`
    }

    body, err := io.ReadAll(r.Body)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    contentType := r.Header.Get("Content-Type")
    if strings.Contains(contentType, "application/json") {
        err = json.Unmarshal(body, &post)
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

        post.Title = values.Get("title")
        post.Content = values.Get("content")
    } else {
        http.Error(w, "Unsupported content type", http.StatusUnsupportedMediaType)
        return
    }

    query := `
        query GetBlog($user_id: uuid!) {
            blogs(where: {user_id: {_eq: $user_id}}) {
                id
            }
        }
    `
    variables := map[string]interface{}{
        "user_id": userID,
    }

    result, err := database.ExecuteGraphQL(query, variables)
    if err != nil {
        log.Printf("Error executing GraphQL query: %v", err)
        http.Error(w, "User does not have a blog", http.StatusBadRequest)
        return
    }

    log.Printf("GraphQL query result: %v", result)

    blogs, ok := result["blogs"].([]interface{})
    if !ok || len(blogs) == 0 {
        log.Printf("No blogs found for user ID %s: %v", userID, result)
        http.Error(w, "User does not have a blog", http.StatusBadRequest)
        return
    }

    blog, ok := blogs[0].(map[string]interface{})
    if !ok {
        log.Printf("Error parsing blog data: %v", blogs[0])
        http.Error(w, "Failed to retrieve blog data", http.StatusInternalServerError)
        return
    }

    blogID, ok := blog["id"].(float64) // Hasura returns numbers as float64
    if !ok {
        log.Printf("Error parsing blog ID: %v", blog["id"])
        http.Error(w, "Failed to retrieve blog ID", http.StatusInternalServerError)
        return
    }

    query = `
        mutation CreatePost($title: String!, $content: String!, $user_id: uuid!, $blog_id: Int!) {
            insert_posts_one(object: {title: $title, content: $content, user_id: $user_id, blog_id: $blog_id}) {
                id
            }
        }
    `
    variables = map[string]interface{}{
        "title":   post.Title,
        "content": post.Content,
        "user_id": userID,
        "blog_id": int(blogID), // Convert float64 to int
    }

    result, err = database.ExecuteGraphQL(query, variables)
    if err != nil {
        log.Printf("Error executing GraphQL mutation: %v", err)
        http.Error(w, "Failed to create post", http.StatusInternalServerError)
        return
    }

    log.Printf("GraphQL mutation result: %v", result)

    postData, ok := result["insert_posts_one"].(map[string]interface{})
    if !ok {
        log.Printf("Error parsing post data from GraphQL result: %v", result)
        http.Error(w, "Failed to create post", http.StatusInternalServerError)
        return
    }

    postID, ok := postData["id"].(float64)
    if !ok {
        log.Printf("Error parsing post ID from GraphQL result: %v", postData)
        http.Error(w, "Failed to create post", http.StatusInternalServerError)
        return
    }

    log.Printf("Successfully created post with ID: %d", int(postID))
    w.Header().Set("HX-Redirect", fmt.Sprintf("/blogs/%d", int(blogID)))
    w.WriteHeader(http.StatusCreated)
}







// PostListHandler handles displaying a list of posts
func PostListHandler(w http.ResponseWriter, r *http.Request) {
    query := `
        query {
            posts {
                id
                title
                content
                user {
                    username
                }
                blog_id
            }
        }
    `

    result, err := database.ExecuteGraphQL(query, nil)
    if err != nil || len(result["errors"].([]interface{})) > 0 {
        http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
        return
    }

    postsData := result["data"].(map[string]interface{})["posts"].([]interface{})
    var posts []models.Post
    for _, postData := range postsData {
        postMap := postData.(map[string]interface{})
        userMap := postMap["user"].(map[string]interface{})

        posts = append(posts, models.Post{
            ID:      uint(postMap["id"].(float64)),
            Title:   postMap["title"].(string),
            Content: postMap["content"].(string),
            User: &models.User{
                Username: userMap["username"].(string),
            },
            BlogID: uint(postMap["blog_id"].(int)),
        })
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

    query := `
        query GetPost($id: Int!) {
            posts_by_pk(id: $id) {
                id
                title
                content
                user {
                    username
                }
                created_at
                blog_id
            }
        }
    `
    variables := map[string]interface{}{
        "id": id,
    }

    result, err := database.ExecuteGraphQL(query, variables)
    if err != nil || len(result["errors"].([]interface{})) > 0 {
        http.Error(w, "Failed to fetch post", http.StatusInternalServerError)
        log.Printf("Post not found: %v", id) // Log the ID not found
        return
    }

    postData := result["data"].(map[string]interface{})["posts_by_pk"].(map[string]interface{})
    userMap := postData["user"].(map[string]interface{})

    post := models.Post{
        ID:       uint(postData["id"].(float64)),
        Title:    postData["title"].(string),
        Content:  postData["content"].(string),
       User: &models.User{
                Username: userMap["username"].(string),
            },
        FormattedCreatedAt: postData["created_at"].(string),
        BlogID:    uint(postData["blog_id"].(int)),
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
