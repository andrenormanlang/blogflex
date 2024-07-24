package handlers

import (
    "encoding/json"
    "net/http"
    "strconv"
    "blogflex/internal/models"
    "github.com/gorilla/mux"
    "blogflex/internal/helpers"
    "time"
    "blogflex/internal/database"
    "fmt"
    "log"
    "strings"
)

func GetCommentsHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    postID, err := strconv.Atoi(vars["postID"])
    if err != nil {
        http.Error(w, "Invalid post ID", http.StatusBadRequest)
        return
    }

    query := `
        query GetComments($post_id: Int!) {
            comments(where: {post_id: {_eq: $post_id}}) {
                id
                content
                post_id
                user_id
            }
        }
    `
    variables := map[string]interface{}{
        "post_id": postID,
    }

    result, err := helpers.GraphQLRequest(query, variables)
    if err != nil || len(result["errors"].([]interface{})) > 0 {
        http.Error(w, "Failed to fetch comments", http.StatusInternalServerError)
        return
    }

    commentsData := result["data"].(map[string]interface{})["comments"].([]interface{})
    var comments []models.Comment
    for _, commentData := range commentsData {
        commentMap := commentData.(map[string]interface{})
        comments = append(comments, models.Comment{
            ID:      uint(commentMap["id"].(float64)),
            Content: commentMap["content"].(string),
            PostID:  uint(commentMap["post_id"].(float64)),
        })
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(comments)
}

func ToggleLikePostHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    postID, err := strconv.Atoi(vars["id"])
    if err != nil {
        helpers.RespondWithError(w, http.StatusBadRequest, "Invalid post ID")
        return
    }

    session, _ := store.Get(r, "session-name")
    userID, ok := session.Values["userID"].(string)
    if !ok {
        helpers.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
        return
    }

    // Check if the user has already liked the post
    query := `
        query CheckLike($post_id: Int!, $user_id: uuid!) {
            likes(where: {post_id: {_eq: $post_id}, user_id: {_eq: $user_id}}) {
                id
            }
        }
    `
    variables := map[string]interface{}{
        "post_id": postID,
        "user_id": userID,
    }

    result, err := database.ExecuteGraphQL(query, variables)
    if err != nil {
        log.Printf("Failed to check like status: %v", err)
        helpers.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve like status")
        return
    }

    log.Printf("CheckLike result: %v", result)

    likes, ok := result["likes"].([]interface{})
    if !ok {
        log.Printf("Unexpected format for 'likes': %v", result)
        helpers.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve like status")
        return
    }

    if len(likes) > 0 {
        // User has liked the post, so unlike it
        likeID := likes[0].(map[string]interface{})["id"].(float64)
        deleteQuery := `
            mutation UnlikePost($id: Int!) {
                delete_likes_by_pk(id: $id) {
                    id
                }
            }
        `
        deleteVariables := map[string]interface{}{
            "id": int(likeID),
        }

        _, err = database.ExecuteGraphQL(deleteQuery, deleteVariables)
        if err != nil {
            log.Printf("Failed to unlike post: %v", err)
            helpers.RespondWithError(w, http.StatusInternalServerError, "Failed to unlike post")
            return
        }
    } else {
        // User has not liked the post, so like it
        likeQuery := `
            mutation LikePost($post_id: Int!, $user_id: uuid!, $created_at: timestamptz!) {
                insert_likes_one(object: {post_id: $post_id, user_id: $user_id, created_at: $created_at}) {
                    id
                }
            }
        `
        likeVariables := map[string]interface{}{
            "post_id":    postID,
            "user_id":    userID,
            "created_at": time.Now().Format(time.RFC3339),
        }

        _, err = database.ExecuteGraphQL(likeQuery, likeVariables)
        if err != nil {
            log.Printf("Failed to like post: %v", err)
            helpers.RespondWithError(w, http.StatusInternalServerError, "Failed to like post")
            return
        }
    }

    // Fetch updated like count
    likesCountQuery := `
        query GetLikesCount($post_id: Int!) {
            posts_with_likes(where: {post_id: {_eq: $post_id}}) {
                likes_count
            }
        }
    `
    likesCountVariables := map[string]interface{}{
        "post_id": postID,
    }

    likesCountResult, err := database.ExecuteGraphQL(likesCountQuery, likesCountVariables)
    if err != nil {
        log.Printf("Failed to fetch likes count: %v", err)
        helpers.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch likes count")
        return
    }

    log.Printf("GetLikesCount result: %v", likesCountResult)

    postsWithLikes, ok := likesCountResult["posts_with_likes"].([]interface{})
    if !ok || len(postsWithLikes) == 0 {
        log.Printf("Unexpected format or empty 'posts_with_likes': %v", likesCountResult)
        helpers.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch likes count")
        return
    }

    likesCount := int(postsWithLikes[0].(map[string]interface{})["likes_count"].(float64))

    // Return the updated HTML for the button
    w.Header().Set("Content-Type", "text/html")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(fmt.Sprintf(`<button id="like-button" hx-post="/protected/posts/%d/like" hx-target="#like-button" hx-swap="outerHTML">
        <i class="fas fa-thumbs-up"></i> <span id="likes-count">%d</span>
    </button>`, postID, likesCount)))
}

func AddCommentHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    postID, err := strconv.Atoi(vars["id"])
    if err != nil {
        helpers.RespondWithError(w, http.StatusBadRequest, "Invalid post ID")
        return
    }

    session, _ := store.Get(r, "session-name")
    userID, ok := session.Values["userID"].(string)
    if !ok {
        helpers.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
        return
    }

    var input struct {
        Content string `json:"content"`
    }

    contentType := r.Header.Get("Content-Type")
    if strings.Contains(contentType, "application/json") {
        if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
            log.Printf("Error decoding JSON input: %v", err)
            http.Error(w, "Invalid input", http.StatusBadRequest)
            return
        }
    } else if strings.Contains(contentType, "application/x-www-form-urlencoded") {
        if err := r.ParseForm(); err != nil {
            log.Printf("Error parsing form input: %v", err)
            http.Error(w, "Invalid input", http.StatusBadRequest)
            return
        }
        input.Content = r.FormValue("content")
    } else {
        log.Printf("Unsupported content type: %v", contentType)
        http.Error(w, "Unsupported content type", http.StatusUnsupportedMediaType)
        return
    }

    comment := models.Comment{
        PostID:    uint(postID),
        UserID:    userID,
        Content:   input.Content,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }

    query := `
        mutation AddComment($content: String!, $post_id: Int!, $user_id: uuid!, $created_at: timestamp!, $updated_at: timestamp!) {
            insert_comments_one(object: {content: $content, post_id: $post_id, user_id: $user_id, created_at: $created_at, updated_at: $updated_at}) {
                id
                content
                user {
                    username
                }
                created_at
                updated_at
            }
        }
    `
    variables := map[string]interface{}{
        "content":    comment.Content,
        "post_id":    comment.PostID,
        "user_id":    comment.UserID,
        "created_at": comment.CreatedAt.Format(time.RFC3339),
        "updated_at": comment.UpdatedAt.Format(time.RFC3339),
    }

    result, err := database.ExecuteGraphQL(query, variables)
    if err != nil {
        log.Printf("Error executing GraphQL request: %v", err)
        helpers.RespondWithError(w, http.StatusInternalServerError, "Failed to add comment")
        return
    }

    newCommentData, ok := result["insert_comments_one"].(map[string]interface{})
    if !ok {
        log.Printf("Unexpected format for insert_comments_one: %v", result)
        helpers.RespondWithError(w, http.StatusInternalServerError, "Failed to add comment")
        return
    }

    comment.ID = uint(newCommentData["id"].(float64))
    comment.CreatedAt, _ = time.Parse(time.RFC3339, newCommentData["created_at"].(string))
    comment.UpdatedAt, _ = time.Parse(time.RFC3339, newCommentData["updated_at"].(string))
    comment.User = &models.User{
        Username: newCommentData["user"].(map[string]interface{})["username"].(string),
    }

    // Return the updated comments in JSON format
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(comment)
}




func GetLikesCountHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    postID, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid post ID", http.StatusBadRequest)
        return
    }

    query := `
        query GetLikesCount($post_id: Int!) {
            likes_count(where: {post_id: {_eq: $post_id}}) {
                aggregate {
                    count
                }
            }
        }
    `
    variables := map[string]interface{}{
        "post_id": postID,
    }

    result, err := helpers.GraphQLRequest(query, variables)
    if err != nil || len(result["errors"].([]interface{})) > 0 {
        http.Error(w, "Failed to fetch likes count", http.StatusInternalServerError)
        return
    }

    count := result["data"].(map[string]interface{})["likes_count"].(map[string]interface{})["aggregate"].(map[string]interface{})["count"].(float64)

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]int{"count": int(count)})
}