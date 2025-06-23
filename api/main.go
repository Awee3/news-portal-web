package main

import (
    "fmt"
    "log"
    "net/http"
)

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello from Go Backend!")
    })

    fmt.Println("Go server running on http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}