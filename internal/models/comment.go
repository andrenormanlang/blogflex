package models

import "gorm.io/gorm"

type Comment struct {
    gorm.Model
    Content string `json:"content"`
    PostID  uint   `json:"post_id"`
    UserID  uint   `json:"user_id"`
}
