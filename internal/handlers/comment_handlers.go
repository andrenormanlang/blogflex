package handlers

import (
    "encoding/json"
    "net/http"
    "strconv"
    "blogflex/internal/models"
    "github.com/gorilla/mux"
    "blogflex/internal/helpers"
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
