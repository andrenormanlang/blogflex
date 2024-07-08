package auth

import "github.com/dgrijalva/jwt-go"

// Define the secret key for signing JWT tokens
var JwtKey = []byte("your-very-secret-key")

// Claims struct to be encoded in the JWT token
type Claims struct {
    Username string `json:"username"`
    jwt.StandardClaims
}
