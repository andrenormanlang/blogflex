package models

import "gorm.io/gorm"

type Comment struct {
    gorm.Model
    Content string `json:"content" gorm:"not null"`
    PostID  uint   `json:"post_id"`
    UserID  uint   `json:"user_id"`
}
