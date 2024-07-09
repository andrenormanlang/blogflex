package auth

import "github.com/dgrijalva/jwt-go"

// Define the secret key for signing JWT tokens
var JwtKey = []byte("your-very-secret-key")


type Claims struct {
    UserID   uint   `json:"user_id"`
    Username string `json:"username"`
    jwt.StandardClaims
}