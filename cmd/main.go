package main

import (
	"fmt"
	"net/http"
)


func main(){

	fmt.Println("hellow from the second image")
    

	http.HandleFunc("/",func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type","application/json")
		w.WriteHeader(200)

		w.Write(make([]byte, 22))
	})

	http.ListenAndServe(":8080",nil)
}