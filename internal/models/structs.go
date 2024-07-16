package models

import (
    "time"
)

type GraphQLRequest struct {
    Query     string                 `json:"query"`
    Variables map[string]interface{} `json:"variables"`
}

type Blog struct {
    ID                 uint      `json:"id"`
    Name               string    `json:"name"`
    Description        string    `json:"description"`
    UserID             uint      `json:"user_id"`
    User               *User     `json:"user"`
    Posts              []Post    `json:"posts"`
    CreatedAt          time.Time `json:"created_at"`
    UpdatedAt          time.Time `json:"updated_at"`
    FormattedCreatedAt string    `json:"-"`
}

type Comment struct {
    ID        uint      `json:"id"`
    Content   string    `json:"content"`
    PostID    uint      `json:"post_id"`
    UserID    string    `json:"user_id"`
    User      *User     `json:"user"`
    Post      *Post     `json:"post"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

type Post struct {
    ID                 uint       `json:"id"`
    Title              string     `json:"title"`
    Content            string     `json:"content"`
    UserID             string     `json:"user_id"`
    User               *User      `json:"user"`
    BlogID             uint       `json:"blog_id"`
    Blog               *Blog      `json:"blog"`
    Comments           []Comment  `json:"comments"`
    CreatedAt          time.Time  `json:"created_at"`
    UpdatedAt          time.Time  `json:"updated_at"`
    FormattedCreatedAt string     `json:"-"`
}

type User struct {
    ID        string    `json:"id"` // UUID is stored as a string
    Username  string    `json:"username"`
    Email     string    `json:"email"`
    Password  string    `json:"password"`
    Blog      *Blog     `json:"blog"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}