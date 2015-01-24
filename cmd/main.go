package main

import (
	"github.com/gophergala/Gift"

	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"time"
)

func storeGeoInCookies(w http.ResponseWriter, r *http.Request) {
	var body string
	if b, err := ioutil.ReadAll(r.Body); err == nil {
		body = string(b)
	}

	u, err := url.Parse("?" + body)
	if err != nil {
		log.Printf("Error parsing setgeo body: %+v", err)
	}
	values := u.Query()

	cookie := http.Cookie{Name: "latitude", Value: values.Get("latitude"), Path: "/"}
	http.SetCookie(w, &cookie)
	cookie = http.Cookie{Name: "longitude", Value: values.Get("longitude"), Path: "/"}
	http.SetCookie(w, &cookie)
	cookie = http.Cookie{Name: "heading", Value: values.Get("heading"), Path: "/"}
	http.SetCookie(w, &cookie)

	log.Printf("Setting geo cookies to: [%s, %s] @ %s", values.Get("latitude"), values.Get("longitude"), values.Get("heading"))

}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	log.Printf("Starting GIFT server")
	counterGiftServer := gift.NewGiftServer(640, 480, &gift.GiftImageCounter{})
	mapGiftServer := gift.NewGiftServer(640, 480, &gift.GiftImageMap{})
	nukeGiftServer := gift.NewGiftServer(480, 360, &gift.GiftImageNuke{})

	http.HandleFunc("/counter.gif", counterGiftServer.Handler)
	http.HandleFunc("/map.gif", mapGiftServer.Handler)
	http.HandleFunc("/nuke.gif", nukeGiftServer.Handler)
	http.HandleFunc("/setgeo", storeGeoInCookies)
	http.Handle("/", http.FileServer(http.Dir("static")))

	http.ListenAndServe(":8080", nil)
}
