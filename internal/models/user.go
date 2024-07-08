package models

import "gorm.io/gorm"

type User struct {
    gorm.Model
    Username string `json:"username"`
    Name     string `json:"name"`
    Email    string `json:"email" gorm:"unique"`
    Password string `json:"-"`
    Posts    []Post `gorm:"foreignKey:UserID"`
}
