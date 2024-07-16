package middleware

import (
    "context"
    "net/http"
    "github.com/dgrijalva/jwt-go"
    "github.com/gorilla/sessions"
    "blogflex/internal/auth"
)

var store = sessions.NewCookieStore([]byte("your-very-secret-key"))

func SessionMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        session, _ := store.Get(r, "session-name")
        ctx := context.WithValue(r.Context(), "session", session)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        session, ok := r.Context().Value("session").(*sessions.Session)
        if !ok {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        tokenStr, ok := session.Values["token"].(string)
        if !ok {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        claims := &auth.Claims{}
        token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
            return auth.JwtKey, nil
        })
        if err != nil || !token.Valid {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        ctx := context.WithValue(r.Context(), "userID", claims.UserID)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
