package main

import (
	"github.com/gophergala/Gift"

	"log"
	"net/http"
)

func main() {
	log.Printf("Starting GIFT server")
	counterGiftServer := gift.NewGiftServer(640, 480, &gift.GiftImageCounter{})
	mapGiftServer := gift.NewGiftServer(640, 480, &gift.GiftImageMap{})

	http.HandleFunc("/counter.gif", counterGiftServer.Handler)
	http.HandleFunc("/map.gif", mapGiftServer.Handler)

	http.ListenAndServe(":8080", nil)
}
