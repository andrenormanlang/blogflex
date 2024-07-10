package models

import (
    "gorm.io/gorm"
)
type Blog struct {
    gorm.Model
    ID                 uint      `gorm:"primaryKey"`
    Name               string    `json:"name"`
    Description        string    `json:"description"`
    UserID             uint      `json:"user_id"`
    User               User      `gorm:"foreignKey:UserID"`
    Posts              []Post
    FormattedCreatedAt string    `gorm:"-"` 
}

type Comment struct {
    gorm.Model
    Content  string `json:"content" gorm:"not null"`
    PostID   uint   `json:"post_id"`
    UserID   uint   `json:"user_id"`
    User     User   `gorm:"foreignKey:UserID"`
    Post     Post   `gorm:"foreignKey:PostID"`
}

type Post struct {
    gorm.Model
    Title    string    `json:"title" gorm:"not null"`
    Content  string    `json:"content" gorm:"not null"`
    UserID   uint      `json:"user_id"`
    User     User      `gorm:"foreignKey:UserID"`
    BlogID   uint      `json:"blog_id"`
    Blog     Blog      `gorm:"foreignKey:BlogID"`
    Comments []Comment `gorm:"foreignKey:PostID"`
}

type User struct {
    gorm.Model
    Username string `json:"username" gorm:"not null"`
    Email    string `json:"email" gorm:"unique;not null"`
    Password string `json:"password" gorm:"not null"`
    Blog     *Blog   `gorm:"foreignKey:UserID" json:"-"` // Use pointer and json:"-"
}
