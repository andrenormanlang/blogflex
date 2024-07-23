package handlers

import (
    "encoding/json"
    "net/http"
    "strconv"
    "blogflex/internal/models"
    "github.com/gorilla/mux"
    "blogflex/internal/helpers"
    "time"
    "log"
    "blogflex/internal/database"
    "fmt"

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

func AddCommentHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    postID, err := strconv.Atoi(vars["postID"])
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

    var comment models.Comment
    if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }

    comment.PostID = uint(postID)
    comment.UserID = userID
    comment.CreatedAt = time.Now()
    comment.UpdatedAt = time.Now()

    query := `
        mutation AddComment($content: String!, $post_id: Int!, $user_id: uuid!, $created_at: timestamptz!, $updated_at: timestamptz!) {
            insert_comments_one(object: {content: $content, post_id: $post_id, user_id: $user_id, created_at: $created_at, updated_at: $updated_at}) {
                id
            }
        }
    `
    variables := map[string]interface{}{
        "content":    comment.Content,
        "post_id":    comment.PostID,
        "user_id":    comment.UserID,
        "created_at": comment.CreatedAt,
        "updated_at": comment.UpdatedAt,
    }

    _, err = helpers.GraphQLRequest(query, variables)
    if err != nil {
        http.Error(w, "Failed to add comment", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
}

func LikePostHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    postID, err := strconv.Atoi(vars["id"])
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

    query := `
        mutation LikePost($post_id: Int!, $user_id: uuid!) {
            insert_likes_one(object: {post_id: $post_id, user_id: $user_id}) {
                id
            }
        }
    `
    variables := map[string]interface{}{
        "post_id": postID,
        "user_id": userID,
    }

    var result map[string]interface{}
    result, err = database.ExecuteGraphQL(query, variables)
    if err != nil {
        log.Printf("Error executing GraphQL mutation: %v", err)
        http.Error(w, "Failed to like post", http.StatusInternalServerError)
        return
    }
    fmt.Printf("GraphQL Result: %+v\n", result)


    response := map[string]string{"status": "success"}
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
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


