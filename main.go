package main

import (
    "log"
    "net/http"
    "blogflex/internal/database"
    "blogflex/internal/router"
)

func main() {
    database.InitHasura()

    r := router.SetupRouter()

    log.Println("Server started at :8080")
    if err := http.ListenAndServe(":8080", r); err != nil {
        log.Fatal("ListenAndServe:", err)
    }
}
