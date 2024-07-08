package middleware

import (
    "net/http"
    "github.com/gorilla/sessions"
)

func AuthMiddleware(store *sessions.CookieStore) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            session, _ := store.Get(r, "session-name")
            userID, ok := session.Values["userID"]
            if !ok || userID == nil {
                http.Error(w, "Unauthorized", http.StatusUnauthorized)
                return
            }
            next.ServeHTTP(w, r)
        })
    }
}
