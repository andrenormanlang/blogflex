package database

import (
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "log"
    "blogflex/internal/models"
)

var DB *gorm.DB

func InitDatabase() {
    dsn := "root:root@tcp(127.0.0.1:3306)/blogflex?charset=utf8mb4&parseTime=True&loc=Local"
    var err error
    DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    // Auto migrate models
    err = DB.AutoMigrate(&models.User{}, &models.Post{}, &models.Comment{})
    if err != nil {
        log.Fatal("Failed to migrate database:", err)
    }
}
