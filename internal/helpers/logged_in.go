package helpers

import (
    "net/http"
    "github.com/dgrijalva/jwt-go"
    "blogflex/internal/auth"
    "github.com/gorilla/sessions"
)

func IsLoggedIn(r *http.Request) bool {
    session, ok := r.Context().Value("session").(*sessions.Session)
    if !ok {
        return false
    }

    tokenStr, ok := session.Values["token"].(string)
    if !ok {
        return false
    }

    claims := &auth.Claims{}
    token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
        return auth.JwtKey, nil
    })
    return err == nil && token.Valid
}
