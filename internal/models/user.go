package models

import "gorm.io/gorm"

type User struct {
    gorm.Model
    Username string `json:"username" gorm:"not null"`
    Email    string `json:"email" gorm:"unique;not null"`
    Password string `json:"password" gorm:"not null"`
    Blog     *Blog  `gorm:"foreignKey:UserID" json:"-"` // Use pointer and json:"-"
}

