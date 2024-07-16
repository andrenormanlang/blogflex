package auth

import (
    "github.com/dgrijalva/jwt-go"
)

var JwtKey = []byte("your-very-secret-key")

type Claims struct {
    UserID   string   `json:"user_id"`
    Username string `json:"username"`
    jwt.StandardClaims
}

type Session struct {
    UserID   uint
    Username string
}
