package main

import (
	"github.com/gophergala/Gift"

	"log"
	"net/http"
)

func main() {
	log.Printf("Starting GIFT server")
	giftServer := gift.NewGiftServer(640, 480, &gift.GiftImageCounter{})

	http.HandleFunc("/", giftServer.Handler)
	http.ListenAndServe(":8080", nil)
}
