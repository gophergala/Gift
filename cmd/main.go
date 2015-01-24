package main

import (
	// "github.com/gophergala/Gift"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Request start")

	log.Printf("Request complete")
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
