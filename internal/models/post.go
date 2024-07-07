package models

import "gorm.io/gorm"

type Post struct {
    gorm.Model
    Title    string    `json:"title"`
    Content  string    `json:"content"`
    UserID   uint      `json:"user_id"`
    User     User      `gorm:"foreignKey:UserID"`
    Comments []Comment `gorm:"foreignKey:PostID"`
}
