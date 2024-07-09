package database

import (
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "log"
)

var DB *gorm.DB

func InitDatabase() *gorm.DB {
    // Replace with your actual DSN
    // Format: username:password@protocol(address)/dbname?param=value
    dsn := "root:root@tcp(127.0.0.1:3306)/blogflex?charset=utf8mb4&parseTime=True&loc=Local"
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    DB = db
    return DB
}
