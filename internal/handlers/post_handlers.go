package handlers

import (
	"blogflex/views"
	"net/http"

	"github.com/a-h/templ"

	// "github.com/gorilla/sessions"
	"blogflex/internal/database"
	"blogflex/internal/helpers"
	"blogflex/internal/models"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
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

	contentType := r.Header.Get("Content-Type")
	if strings.Contains(contentType, "application/json") {
		err := json.NewDecoder(r.Body).Decode(&post)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	} else if strings.Contains(contentType, "application/x-www-form-urlencoded") {
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		post.Title = r.FormValue("title")
		post.Content = r.FormValue("content")

		// Ensure content is properly decoded
		if decodedContent, err := url.QueryUnescape(post.Content); err == nil {
			post.Content = decodedContent
		}
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

func UploadImageHandler(w http.ResponseWriter, r *http.Request) {
	// Limit the size of the uploaded file
	r.ParseMultipartForm(10 << 20) // 10 MB

	file, handler, err := r.FormFile("upload")
	if err != nil {
		http.Error(w, "Unable to upload file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Create a unique file name
	fileName := filepath.Join("uploads", handler.Filename)
	out, err := os.Create(fileName)
	if err != nil {
		http.Error(w, "Unable to create the file for writing. Check your write access privilege", http.StatusInternalServerError)
		return
	}
	defer out.Close()

	// Write the content from the file to the new file
	_, err = io.Copy(out, file)
	if err != nil {
		http.Error(w, "Unable to write file", http.StatusInternalServerError)
		return
	}

	// Return the URL of the uploaded file
	response := map[string]interface{}{
		"url": "/uploads/" + handler.Filename,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func EditPostFormHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	query := `
        query GetPost($id: Int!) {
            posts_by_pk(id: $id) {
                id
                title
                content
                user {
                    id
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
	if err != nil {
		log.Printf("Failed to execute GraphQL query: %v", err)
		http.Error(w, "Failed to fetch post", http.StatusInternalServerError)
		return
	}

	// Check if the result directly contains the required data
	postData, ok := result["posts_by_pk"].(map[string]interface{})
	if (!ok) || postData == nil {
		log.Printf("posts_by_pk is nil or not a map: %v", result)
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	userMap, ok := postData["user"].(map[string]interface{})
	if (!ok) || userMap == nil {
		log.Printf("user is nil or not a map: %v", postData["user"])
		http.Error(w, "User data not found", http.StatusInternalServerError)
		return
	}

	// Ensure only the owner can edit the post
	session, _ := store.Get(r, "session-name")
	loggedInUserID, ok := session.Values["userID"].(string)
	if (!ok) || loggedInUserID != userMap["id"].(string) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	post := models.Post{
		ID:      uint(postData["id"].(float64)),
		Title:   postData["title"].(string),
		Content: postData["content"].(string),
	}

	component := views.EditPost(post)
	templ.Handler(component).ServeHTTP(w, r)
}

func EditPostHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

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

	// Log the request details for debugging
	log.Printf("Edit Post Request - Method: %s, Content-Type: %s", r.Method, r.Header.Get("Content-Type"))

	contentType := r.Header.Get("Content-Type")
	if strings.Contains(contentType, "application/json") {
		err := json.NewDecoder(r.Body).Decode(&post)
		if err != nil {
			log.Printf("Error decoding JSON: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	} else if strings.Contains(contentType, "application/x-www-form-urlencoded") {
		if err := r.ParseForm(); err != nil {
			log.Printf("Error parsing form: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		post.Title = r.FormValue("title")
		post.Content = r.FormValue("content")
	} else if strings.Contains(contentType, "multipart/form-data") {
		// Handle multipart form data (which is what FormData sends)
		if err := r.ParseMultipartForm(10 << 20); err != nil { // 10 MB max memory
			log.Printf("Error parsing multipart form: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		post.Title = r.FormValue("title")
		post.Content = r.FormValue("content")
	} else {
		log.Printf("Unsupported content type: %s", contentType)
		http.Error(w, "Unsupported content type", http.StatusUnsupportedMediaType)
		return
	}

	// Log the extracted data
	log.Printf("Post data - Title: %s, Content length: %d", post.Title, len(post.Content))

	query := `
        mutation UpdatePost($id: Int!, $title: String!, $content: String!, $user_id: uuid!) {
            update_posts(where: {id: {_eq: $id}, user_id: {_eq: $user_id}}, _set: {title: $title, content: $content}) {
                affected_rows
            }
        }
    `
	variables := map[string]interface{}{
		"id":      id,
		"title":   post.Title,
		"content": post.Content,
		"user_id": userID,
	}

	result, err := database.ExecuteGraphQL(query, variables)
	if err != nil {
		log.Printf("Error executing GraphQL mutation: %v", err)
		http.Error(w, "Failed to update post", http.StatusInternalServerError)
		return
	}

	// Check affected rows
	if updateResult, ok := result["update_posts"].(map[string]interface{}); ok {
		if affectedRows, ok := updateResult["affected_rows"].(float64); ok {
			log.Printf("Update affected %d rows", int(affectedRows))
			if affectedRows == 0 {
				log.Printf("No rows were updated. This could be because the post doesn't exist or the user doesn't have permission.")
				http.Error(w, "No changes were made. The post may not exist or you may not have permission to edit it.", http.StatusNotFound)
				return
			}
		}
	}

	log.Printf("GraphQL mutation result: %v", result)
	w.Header().Set("HX-Redirect", fmt.Sprintf("/posts/%d", id)) // Redirect to post detail page
	w.WriteHeader(http.StatusOK)
}

func DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	session, _ := store.Get(r, "session-name")
	userID, ok := session.Values["userID"].(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Fetch blog ID before deleting the post
	query := `
        query GetPost($id: Int!) {
            posts_by_pk(id: $id) {
                blog_id
            }
        }
    `
	variables := map[string]interface{}{
		"id": id,
	}

	result, err := database.ExecuteGraphQL(query, variables)
	if err != nil {
		log.Printf("Error executing GraphQL query: %v", err)
		http.Error(w, "Failed to fetch post", http.StatusInternalServerError)
		return
	}

	postData, ok := result["posts_by_pk"].(map[string]interface{})
	if !ok || postData == nil {
		log.Printf("posts_by_pk is nil or not a map: %v", result)
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	blogID, ok := postData["blog_id"].(float64)
	if !ok {
		log.Printf("Error parsing blog ID from GraphQL result: %v", postData)
		http.Error(w, "Failed to retrieve blog ID", http.StatusInternalServerError)
		return
	}

	// Delete the post
	query = `
        mutation DeletePost($id: Int!, $user_id: uuid!) {
            delete_posts(where: {id: {_eq: $id}, user_id: {_eq: $user_id}}) {
                affected_rows
            }
        }
    `
	variables = map[string]interface{}{
		"id":      id,
		"user_id": userID,
	}

	result, err = database.ExecuteGraphQL(query, variables)
	if err != nil {
		log.Printf("Error executing GraphQL mutation: %v", err)
		http.Error(w, "Failed to delete post", http.StatusInternalServerError)
		return
	}

	log.Printf("GraphQL mutation result: %v", result)
	w.Header().Set("HX-Redirect", fmt.Sprintf("/blogs/%d", int(blogID))) // Redirect to blog page
	w.WriteHeader(http.StatusOK)
}

func PostDetailHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		log.Printf("Invalid post ID: %v", vars["id"]) // Log the invalid ID
		return
	}

	// Query to get post details
	postQuery := `
        query GetPost($id: Int!) {
            posts_by_pk(id: $id) {
                id
                title
                content
                user {
                    id
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

	postResult, err := database.ExecuteGraphQL(postQuery, variables)
	if err != nil {
		log.Printf("Failed to execute GraphQL query: %v", err)
		http.Error(w, "Failed to fetch post", http.StatusInternalServerError)
		return
	}

	postData, ok := postResult["posts_by_pk"].(map[string]interface{})
	if !ok || postData == nil {
		log.Printf("posts_by_pk is nil or not a map: %v", postResult)
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	userMap, ok := postData["user"].(map[string]interface{})
	if !ok || userMap == nil {
		log.Printf("user is nil or not a map: %v", postData["user"])
		http.Error(w, "User data not found", http.StatusInternalServerError)
		return
	}

	createdAtStr, ok := postData["created_at"].(string)
	var formattedCreatedAt string
	if ok {
		createdAt, err := time.Parse("2006-01-02T15:04:05", createdAtStr)
		if err != nil {
			log.Printf("Error parsing post created_at time: %v", err)
		} else {
			formattedCreatedAt = helpers.FormatTime(createdAt)
		}
	}

	post := models.Post{
		ID:      uint(postData["id"].(float64)),
		Title:   postData["title"].(string),
		Content: postData["content"].(string),
		User: &models.User{
			ID:       userMap["id"].(string),
			Username: userMap["username"].(string),
		},
		FormattedCreatedAt: formattedCreatedAt,
		BlogID:             uint(postData["blog_id"].(float64)),
	}

	// Fetch likes count for the post
	likesCountQuery := `
		query GetLikesCount($post_id: Int!) {
			likes_aggregate(where: {post_id: {_eq: $post_id}}) {
				aggregate {
					count
				}
			}
		}
	`
	likesCountVariables := map[string]interface{}{
		"post_id": id,
	}

	likesCountResult, err := database.ExecuteGraphQL(likesCountQuery, likesCountVariables)
	if err != nil {
		log.Printf("Failed to fetch likes count: %v", err)
		// Continue without likes count
	} else {
		if likesAggregate, ok := likesCountResult["likes_aggregate"].(map[string]interface{}); ok {
			if aggregate, ok := likesAggregate["aggregate"].(map[string]interface{}); ok {
				if count, ok := aggregate["count"].(float64); ok {
					post.LikesCount = int(count)
				}
			}
		}
	}

	// Check if the current user has liked this post
	session, _ := store.Get(r, "session-name")
	loggedInUserID, ok := session.Values["userID"].(string)
	isOwner := ok && loggedInUserID == post.User.ID

	// Only check if user has liked the post if they're logged in and not the owner
	userHasLiked := false
	if ok && !isOwner {
		hasLikedQuery := `
			query CheckUserLike($post_id: Int!, $user_id: uuid!) {
				likes(where: {post_id: {_eq: $post_id}, user_id: {_eq: $user_id}}) {
					id
				}
			}
		`
		hasLikedVariables := map[string]interface{}{
			"post_id": id,
			"user_id": loggedInUserID,
		}

		hasLikedResult, err := database.ExecuteGraphQL(hasLikedQuery, hasLikedVariables)
		if err == nil {
			if likes, ok := hasLikedResult["likes"].([]interface{}); ok {
				userHasLiked = len(likes) > 0
			}
		}
	}
	post.UserHasLiked = userHasLiked

	// Fetch comments for the post
	commentsQuery := `
        query GetComments($post_id: Int!) {
            comments(where: {post_id: {_eq: $post_id}}) {
                id
                content
                user {
                    id
                    username
                }
                created_at
                updated_at
                user_id
            }
        }
    `
	commentsVariables := map[string]interface{}{
		"post_id": id,
	}

	commentsResult, err := database.ExecuteGraphQL(commentsQuery, commentsVariables)
	if err != nil {
		log.Printf("Failed to execute GraphQL query: %v", err)
		http.Error(w, "Failed to fetch comments", http.StatusInternalServerError)
		return
	}

	commentsData, ok := commentsResult["comments"].([]interface{})
	if !ok {
		log.Printf("Unexpected format for comments: %v", commentsResult)
		http.Error(w, "Failed to fetch comments", http.StatusInternalServerError)
		return
	}

	var comments []models.Comment
	for _, c := range commentsData {
		commentMap := c.(map[string]interface{})
		createdAtStr, ok := commentMap["created_at"].(string)
		var formattedCommentCreatedAt string
		if ok {
			// The GraphQL API returns timestamps in RFC3339 format
			createdAt, err := time.Parse(time.RFC3339, createdAtStr)
			if err != nil {
				// Try alternative format if the first one fails
				createdAt, err = time.Parse("2006-01-02T15:04:05Z", createdAtStr)
				if err != nil {
					// Try one more format
					createdAt, err = time.Parse("2006-01-02T15:04:05", createdAtStr)
					if err != nil {
						log.Printf("Error parsing comment created_at time: %v", err)
						log.Printf("Problematic timestamp: %s", createdAtStr)
					} else {
						formattedCommentCreatedAt = helpers.FormatTime(createdAt)
					}
				} else {
					formattedCommentCreatedAt = helpers.FormatTime(createdAt)
				}
			} else {
				formattedCommentCreatedAt = helpers.FormatTime(createdAt)
			}
		}

		comments = append(comments, models.Comment{
			ID:      uint(commentMap["id"].(float64)),
			Content: commentMap["content"].(string),
			User: &models.User{
				ID:       commentMap["user"].(map[string]interface{})["id"].(string),
				Username: commentMap["user"].(map[string]interface{})["username"].(string),
			},
			UserID:             commentMap["user_id"].(string),
			FormattedCreatedAt: formattedCommentCreatedAt,
		})
	}
	post.Comments = comments

	// Check if the logged-in user is the owner of the post
	session, _ = store.Get(r, "session-name")
	loggedInUserID, ok = session.Values["userID"].(string)
	isOwner = ok && loggedInUserID == post.User.ID

	log.Printf("Post found: ID=%d, Title=%s", post.ID, post.Title) // Log the found post

	// Render HTML template
	loggedIn := helpers.IsLoggedIn(r)
	component := views.PostDetail(post, loggedIn, isOwner, loggedInUserID)
	templ.Handler(component).ServeHTTP(w, r)
}
