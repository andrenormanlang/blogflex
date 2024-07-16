
package models

import (
    "time"
)

type Blog struct {
    ID                 uint      `json:"id"`
    Name               string    `json:"name"`
    Description        string    `json:"description"`
    UserID             uint      `json:"user_id"`
    User               *User     `json:"user"` // Use pointer to avoid recursive type
    Posts              []Post    `json:"posts"`
    CreatedAt          time.Time `json:"created_at"`
    UpdatedAt          time.Time `json:"updated_at"`
    FormattedCreatedAt string    `json:"-"` // Omits this field from JSON serialization
}

type Comment struct {
    ID        uint      `json:"id"`
    Content   string    `json:"content"`
    PostID    uint      `json:"post_id"`
    UserID    uint      `json:"user_id"`
    User      *User     `json:"user"` // Use pointer to avoid recursive type
    Post      *Post     `json:"post"` // Use pointer to avoid recursive type
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

type Post struct {
    ID                 uint       `json:"id"`
    Title              string     `json:"title"`
    Content            string     `json:"content"`
    UserID             uint       `json:"user_id"`
    User               *User      `json:"user"` // Use pointer to avoid recursive type
    BlogID             uint       `json:"blog_id"`
    Blog               *Blog      `json:"blog"` // Use pointer to avoid recursive type
    Comments           []Comment  `json:"comments"`
    CreatedAt          time.Time  `json:"created_at"`
    UpdatedAt          time.Time  `json:"updated_at"`
    FormattedCreatedAt string     `json:"-"` // Omits this field from JSON serialization
}

type User struct {
    ID        uint       `json:"id"`
    Username  string     `json:"username"`
    Email     string     `json:"email"`
    Password  string     `json:"password"`
    Blog      *Blog      `json:"blog"` // Use pointer to avoid recursive type
    CreatedAt time.Time  `json:"created_at"`
    UpdatedAt time.Time  `json:"updated_at"`
}
