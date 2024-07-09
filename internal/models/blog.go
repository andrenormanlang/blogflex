package models

import "gorm.io/gorm"

type Blog struct {
    gorm.Model
    Name              string `json:"name"`
    Description       string `json:"description"`
    UserID            uint   `json:"user_id"`
    User              *User  `gorm:"foreignKey:UserID" json:"-"` // Use pointer and json:"-"
}
