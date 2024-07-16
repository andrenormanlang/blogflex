// package handlers

// import (
//     "net/http"
//     "github.com/a-h/templ"
//     "blogflex/internal/database"
//     "blogflex/internal/models"
//     "encoding/json"
//     "github.com/gorilla/mux"
//     "strconv"
//     "strings"
//     "log"
//     "blogflex/views"
//     "io"
//     "net/url"
//     "github.com/gorilla/sessions"
// )

// // CreatePostFormHandler handles the form submission for creating a post
// func CreatePostFormHandler(w http.ResponseWriter, r *http.Request) {
//     session := r.Context().Value("session").(*sessions.Session)
//     userID, ok := session.Values["userID"].(uint)
//     if !ok {
//         http.Error(w, "Unauthorized", http.StatusUnauthorized)
//         return
//     }

//     component := views.CreatePost(userID)
//     templ.Handler(component).ServeHTTP(w, r)
// }

// // CreatePostHandler handles creating a new post
// func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
//     session := r.Context().Value("session").(*sessions.Session)
//     userID, ok := session.Values["userID"].(uint)
//     if !ok {
//         http.Error(w, "Unauthorized", http.StatusUnauthorized)
//         return
//     }

//     var post models.Post

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
//         err = json.Unmarshal(body, &post)
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

//         post.Title = values.Get("title")
//         post.Content = values.Get("content")
//     } else {
//         http.Error(w, "Unsupported content type", http.StatusUnsupportedMediaType)
//         return
//     }

//     post.UserID = userID // Set the user ID from session

//     // Ensure the post is linked to the user's blog
//     var blog models.Blog
//     if err := database.DB.Where("user_id = ?", userID).First(&blog).Error; err != nil {
//         http.Error(w, "User does not have a blog", http.StatusBadRequest)
//         return
//     }
//     post.BlogID = blog.ID

//     result := database.DB.Create(&post)
//     if result.Error != nil {
//         http.Error(w, result.Error.Error(), http.StatusInternalServerError)
//         return
//     }

//     // Respond with a redirect
//     w.Header().Set("HX-Redirect", "/blogs/"+strconv.Itoa(int(blog.ID)))
//     w.WriteHeader(http.StatusCreated)
// }

// // PostListHandler handles displaying a list of posts
// func PostListHandler(w http.ResponseWriter, r *http.Request) {
//     var posts []models.Post
//     result := database.DB.Find(&posts)
//     if result.Error != nil {
//         http.Error(w, result.Error.Error(), http.StatusInternalServerError)
//         return
//     }

//     // Log the posts to ensure they have IDs
//     for _, post := range posts {
//         log.Printf("Post ID: %d, Title: %s", post.ID, post.Title)
//     }

//     // Check the Accept header to determine the response format
//     acceptHeader := r.Header.Get("Accept")
//     log.Printf("Accept header: %s", acceptHeader) // Log the Accept header value

//     if strings.Contains(acceptHeader, "application/json") {
//         log.Println("Returning JSON response") // Log the branch being executed
//         // Set the content type to JSON
//         w.Header().Set("Content-Type", "application/json")
//         // Encode the posts as JSON and write to the response
//         if err := json.NewEncoder(w).Encode(posts); err != nil {
//             http.Error(w, err.Error(), http.StatusInternalServerError)
//         }
//     } else {
//         log.Println("Returning HTML response") // Log the branch being executed
//         // Render HTML template
//         component := views.PostList(posts) // Correctly refer to the templates.PostList component
//         templ.Handler(component).ServeHTTP(w, r)
//     }
// }

// // PostDetailHandler handles displaying the details of a single post
// func PostDetailHandler(w http.ResponseWriter, r *http.Request) {
//     vars := mux.Vars(r)
//     id, err := strconv.Atoi(vars["id"])
//     if err != nil {
//         http.Error(w, "Invalid post ID", http.StatusBadRequest)
//         log.Printf("Invalid post ID: %v", vars["id"]) // Log the invalid ID
//         return
//     }

//     var post models.Post
//     result := database.DB.First(&post, id)
//     if result.Error != nil {
//         http.Error(w, result.Error.Error(), http.StatusNotFound)
//         log.Printf("Post not found: %v", id) // Log the ID not found
//         return
//     }

//     log.Printf("Post found: ID=%d, Title=%s", post.ID, post.Title) // Log the found post

//     // Check the Accept header to determine the response format
//     acceptHeader := r.Header.Get("Accept")
//     log.Printf("Accept header: %s", acceptHeader) // Log the Accept header value

//     if strings.Contains(acceptHeader, "application/json") {
//         log.Println("Returning JSON response") // Log the branch being executed
//         // Set the content type to JSON
//         w.Header().Set("Content-Type", "application/json")
//         // Encode the post as JSON and write to the response
//         if err := json.NewEncoder(w).Encode(post); err != nil {
//             http.Error(w, err.Error(), http.StatusInternalServerError)
//         }
//     } else {
//         log.Println("Returning HTML response") // Log the branch being executed
//         // Render HTML template
//         component := views.PostDetail(post) // Correctly refer to the templates.PostDetail component
//         templ.Handler(component).ServeHTTP(w, r)
//     }
// }

package handlers

import (
    "net/http"
    "github.com/a-h/templ"
    "blogflex/views"
    "github.com/gorilla/sessions"
    "encoding/json"
    "io"
    "net/url"
    "strings"
    "log"
    "blogflex/internal/models"
    "strconv"
    "github.com/gorilla/mux"
    "blogflex/internal/database"
)

// CreatePostFormHandler handles the form submission for creating a post
func CreatePostFormHandler(w http.ResponseWriter, r *http.Request) {
    session := r.Context().Value("session").(*sessions.Session)
    userIDStr, ok := session.Values["userID"].(string)
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    userID, err := strconv.ParseUint(userIDStr, 10, 32)
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusInternalServerError)
        return
    }

    component := views.CreatePost(uint(userID))
    templ.Handler(component).ServeHTTP(w, r)
}

// CreatePostHandler handles creating a new post
func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
    session := r.Context().Value("session").(*sessions.Session)
    userID, ok := session.Values["userID"].(string)
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    var post struct {
        Title   string `json:"title"`
        Content string `json:"content"`
    }

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
        err = json.Unmarshal(body, &post)
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

        post.Title = values.Get("title")
        post.Content = values.Get("content")
    } else {
        http.Error(w, "Unsupported content type", http.StatusUnsupportedMediaType)
        return
    }

    // Ensure the post is linked to the user's blog
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
    if err != nil || len(result["errors"].([]interface{})) > 0 {
        http.Error(w, "User does not have a blog", http.StatusBadRequest)
        return
    }

    blogs := result["data"].(map[string]interface{})["blogs"].([]interface{})
    if len(blogs) == 0 {
        http.Error(w, "User does not have a blog", http.StatusBadRequest)
        return
    }

    blogID := blogs[0].(map[string]interface{})["id"].(string)

    // Insert the post
    query = `
        mutation CreatePost($title: String!, $content: String!, $user_id: uuid!, $blog_id: uuid!) {
            insert_posts_one(object: {title: $title, content: $content, user_id: $user_id, blog_id: $blog_id}) {
                id
            }
        }
    `
    variables = map[string]interface{}{
        "title":   post.Title,
        "content": post.Content,
        "user_id": userID,
        "blog_id": blogID,
    }

    result, err = database.ExecuteGraphQL(query, variables)
    if err != nil || len(result["errors"].([]interface{})) > 0 {
        http.Error(w, "Failed to create post", http.StatusInternalServerError)
        return
    }

    // Respond with a redirect
    w.Header().Set("HX-Redirect", "/blogs/"+blogID)
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
            BlogID: uint(postMap["blog_id"].(float64)),
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
        BlogID:    uint(postData["blog_id"].(float64)),
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






// CreatePostHandler handles creating a new post earlier
// func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
//     var post models.Post

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
//         err = json.Unmarshal(body, &post)
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

//         post.Title = values.Get("title")
//         post.Content = values.Get("content")
//         post.UserID = 1 // Hardcoded user ID for demonstration
//     } else {
//         http.Error(w, "Unsupported content type", http.StatusUnsupportedMediaType)
//         return
//     }

//     result := database.DB.Create(&post)
//     if result.Error != nil {
//         http.Error(w, result.Error.Error(), http.StatusInternalServerError)
//         return
//     }

//     // Respond with a success message
//     w.Header().Set("Content-Type", "text/html")
//     w.WriteHeader(http.StatusCreated)
//     response := `<div class="bg-green-100 border-t border-b border-green-500 text-green-700 px-4 py-3" role="alert">
//                     <p class="font-bold">Success!</p>
//                     <p class="text-sm">Post created successfully.</p>
//                  </div>`
//     w.Write([]byte(response))
// }