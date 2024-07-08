package models

import "gorm.io/gorm"

type Post struct {
    gorm.Model
     ID      uint   `gorm:"primaryKey"`
    Title    string    `json:"title"`
    Content  string    `json:"content"`
    UserID   uint      `json:"user_id"`
    User     User      `gorm:"foreignKey:UserID"`
    Comments []Comment `gorm:"foreignKey:PostID"`
}
