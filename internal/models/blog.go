package models

import (
    "gorm.io/gorm"
)

type Blog struct {
    gorm.Model
    Name              string `json:"name"`
    Description       string `json:"description"`
    UserID            uint   `json:"user_id"`
    User              User   `gorm:"foreignKey:UserID"`
    FormattedCreatedAt string `gorm:"-"` // This is a computed field, so we use `gorm:"-"`
}
