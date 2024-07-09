package models

import "gorm.io/gorm"

type Post struct {
    gorm.Model
    Title    string `json:"title" gorm:"not null"`
    Content  string `json:"content" gorm:"not null"`
    UserID   uint   `json:"user_id"`
    User     User   `gorm:"foreignKey:UserID"`
    BlogID   uint   `json:"blog_id"`
    Blog     Blog   `gorm:"foreignKey:BlogID"`
    Comments []Comment `gorm:"foreignKey:PostID"`
}
