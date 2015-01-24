package gift

import (
	"image"
	"image/color/palette"
	"image/gif"
	"log"
	"net"
	"net/http"
)

func (gs *GiftServer) Handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Request start")
	gs.source.Setup(gs.width, gs.height)
	w.Header().Set("Content-Type", "image/gif")

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		log.Printf("Error: %+v", err)
	}
	if host == "::1" {
		host = "8.8.8.8"
	}
	ip := net.ParseIP(host)

	record, err := gs.db.City(ip)
	if err != nil {
		log.Printf("Error looking up ip: %+v", err)
	}
	gs.source.Geo(record)

	g := gif.GIF{}
	g.Image = append(g.Image, image.NewPaletted(image.Rect(0, 0, gs.width, gs.height), palette.Plan9))
	g.Delay = append(g.Delay, 100)
	g.LoopCount = 0

	images := make(chan *image.Paletted)

	go gs.source.Pipe(images)

	err = EncodeAll(w, &g, images)
	if err != nil {
		log.Printf("Err: %+v", err)
	}
	log.Printf("Request complete")

}
