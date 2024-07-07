package models

import "gorm.io/gorm"

type Comment struct {
    gorm.Model
    Content string `json:"content"`
    PostID  uint
    UserID  uint
}
