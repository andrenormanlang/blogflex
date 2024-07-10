package middleware

import (
    "context"
    "net/http"
    "blogflex/internal/auth"
    "log"
    "github.com/dgrijalva/jwt-go"
)

var JwtKey = []byte("your-very-secret-key")

func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        cookie, err := r.Cookie("token")
        if err != nil {
            if err == http.ErrNoCookie {
                http.Error(w, "Unauthorized", http.StatusUnauthorized)
                return
            }
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        tokenStr := cookie.Value
        claims := &auth.Claims{}

        token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
            return JwtKey, nil
        })

        if err != nil {
            if err == jwt.ErrSignatureInvalid {
                http.Error(w, "Unauthorized", http.StatusUnauthorized)
                return
            }
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        if !token.Valid {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        log.Printf("Authenticated user: %s", claims.Username)

        session := &auth.Session{
            UserID:   claims.UserID,
            Username: claims.Username,
        }
        ctx := context.WithValue(r.Context(), "session", session)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}


// package middleware

// import (
//     "net/http"
//     "github.com/gorilla/sessions"
// )

// func AuthMiddleware(store *.CookieStore) func(http.Handler) http.Handler {
//     return func(next http.Handler) http.Handler {
//         return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//             session, _ := store.Get(r, "session-name")
//             userID, ok := session.Values["userID"]
//             if !ok || userID == nil {
//                 http.Error(w, "Unauthorized", http.StatusUnauthorized)
//                 return
//             }
//             next.ServeHTTP(w, r)
//         })
//     }
// }
