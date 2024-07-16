// package handlers

// import (
//     "encoding/json"
//     "net/http"
//     "blogflex/internal/database"
//     "github.com/gorilla/mux"
//     "blogflex/internal/models"
//     "strconv"
// )

// func CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
//     var comment models.Comment
//     err := json.NewDecoder(r.Body).Decode(&comment)
//     if err != nil {
//         http.Error(w, err.Error(), http.StatusBadRequest)
//         return
//     }

//     // Ensure the post ID is valid
//     var post models.Post
//     if err := database.DB.First(&post, comment.PostID).Error; err != nil {
//         http.Error(w, "Invalid post ID", http.StatusBadRequest)
//         return
//     }

//     // Ensure the user ID is valid
//     var user models.User
//     if err := database.DB.First(&user, comment.UserID).Error; err != nil {
//         http.Error(w, "Invalid user ID", http.StatusBadRequest)
//         return
//     }

//     result := database.DB.Create(&comment)
//     if result.Error != nil {
//         http.Error(w, result.Error.Error(), http.StatusInternalServerError)
//         return
//     }

//     w.WriteHeader(http.StatusCreated)
//     json.NewEncoder(w).Encode(comment)
// }

// func GetCommentsHandler(w http.ResponseWriter, r *http.Request) {
//     vars := mux.Vars(r)
//     postID, err := strconv.Atoi(vars["postID"])
//     if err != nil {
//         http.Error(w, "Invalid post ID", http.StatusBadRequest)
//         return
//     }

//     var comments []models.Comment
//     if err := database.DB.Where("post_id = ?", postID).Find(&comments).Error; err != nil {
//         http.Error(w, err.Error(), http.StatusInternalServerError)
//         return
//     }

//     w.Header().Set("Content-Type", "application/json")
//     json.NewEncoder(w).Encode(comments)
// }

package handlers

import (
    "encoding/json"
    "net/http"
    "strconv"
    "blogflex/internal/models"
    "github.com/gorilla/mux"
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

    result, err := graphqlRequest(query, variables)
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
            UserID:  uint(commentMap["user_id"].(float64)),
        })
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(comments)
}