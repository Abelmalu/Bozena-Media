package main

import (
    "fmt"
    "net/http"
)

// 1. Define a function with the correct signature
func helloHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "Hello, world!")
}

func main() {
    // 2. Use http.HandleFunc (it calls HandlerFunc internally)
    http.HandleFunc("/hello", helloHandler)
    
    http.ListenAndServe(":8080", nil)
}
type HandlerFunc func(w http.ResponseWriter, r *http.Request)


//http.HanlderFunc is used to map a path to a hanlder function 