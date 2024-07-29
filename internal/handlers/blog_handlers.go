package handlers

import (
	"blogflex/internal/auth"
	"blogflex/internal/helpers"
	"blogflex/internal/models"
	"blogflex/views"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/a-h/templ"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
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
// BlogPageHandler handles displaying a single blog page with its posts
func BlogPageHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    blogID, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid blog ID", http.StatusBadRequest)
        return
    }

    // Query to get blog details and posts
    blogQuery := `
        query GetBlogWithPosts($id: Int!) {
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
                    created_at
                    comments_aggregate {
                        aggregate {
                            count
                        }
                    }
                }
            }
        }
    `
    blogVariables := map[string]interface{}{
        "id": blogID,
    }

    blogResult, err := helpers.GraphQLRequest(blogQuery, blogVariables)
    if err != nil {
        log.Printf("Failed to fetch blog: %v", err)
        http.Error(w, "Failed to fetch blog", http.StatusInternalServerError)
        return
    }

    log.Printf("GraphQL blog result: %+v\n", blogResult)

    if blogResult["data"] == nil {
        log.Println("blogResult[data] is nil")
        http.Error(w, "Failed to fetch blog data", http.StatusInternalServerError)
        return
    }

    blogData, ok := blogResult["data"].(map[string]interface{})["blogs_by_pk"].(map[string]interface{})
    if !ok || blogData == nil {
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

    // Collect post IDs to fetch likes count
    var postIDs []int
    for _, postData := range postsData {
        postMap := postData.(map[string]interface{})
        postID := int(postMap["id"].(float64))
        postIDs = append(postIDs, postID)
    }

    var likesMap map[int]int
    if len(postIDs) > 0 {
        // Query to get likes count for the posts
        likesQuery := `
            query GetLikesCounts($post_ids: [Int!]!) {
                posts_with_likes(where: {post_id: {_in: $post_ids}}) {
                    post_id
                    likes_count
                }
            }
        `
        likesVariables := map[string]interface{}{
            "post_ids": postIDs,
        }

        likesResult, err := helpers.GraphQLRequest(likesQuery, likesVariables)
        if err != nil {
            log.Printf("Failed to fetch likes count: %v", err)
            http.Error(w, "Failed to fetch likes count", http.StatusInternalServerError)
            return
        }

        log.Printf("GraphQL likes result: %+v\n", likesResult)

        likesData, ok := likesResult["data"].(map[string]interface{})["posts_with_likes"].([]interface{})
        if !ok || likesData == nil {
            log.Printf("posts_with_likes is nil or not a slice: %v", likesResult["data"])
            http.Error(w, "Failed to fetch likes data", http.StatusInternalServerError)
            return
        }

        // Create a map of post IDs to likes count
        likesMap = make(map[int]int)
        for _, like := range likesData {
            likeMap := like.(map[string]interface{})
            postID := int(likeMap["post_id"].(float64))
            likesCount := 0
            if likeMap["likes_count"] != nil {
                likesCount = int(likeMap["likes_count"].(float64))
            }
            likesMap[postID] = likesCount
        }
    } else {
        likesMap = make(map[int]int)
    }

    // Process posts and combine likes count
    var posts []models.Post
    for _, postData := range postsData {
        postMap, ok := postData.(map[string]interface{})
        if !ok {
            log.Printf("postData is not a map: %v\n", postData)
            continue
        }

        postID := int(postMap["id"].(float64))

        createdAtStr, ok := postMap["created_at"].(string)
        var formattedCreatedAt string
        if ok {
            createdAt, err := time.Parse("2006-01-02T15:04:05", createdAtStr)
            if err != nil {
                log.Printf("Error parsing post created_at time: %v", err)
            } else {
                formattedCreatedAt = helpers.FormatTime(createdAt)
            }
        }

        posts = append(posts, models.Post{
            ID:                uint(postID),
            Title:             postMap["title"].(string),
            Content:           postMap["content"].(string),
            FormattedCreatedAt: formattedCreatedAt,
            CommentsCount:     int(postMap["comments_aggregate"].(map[string]interface{})["aggregate"].(map[string]interface{})["count"].(float64)),
            LikesCount:        likesMap[postID], // Get likes count from the map
        })
    }

    blog := models.Blog{
        ID:          uint(blogData["id"].(float64)),
        Name:        blogData["name"].(string),
        Description: blogData["description"].(string),
        User: &models.User{
            ID:       userMap["id"].(string),
            Username: userMap["username"].(string),
        },
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


