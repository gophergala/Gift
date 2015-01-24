package gift

import (
	"github.com/oschwald/geoip2-golang"

	"image"
	"log"
)

type GiftImageSource interface {
	Setup(width, height int)
	Geo(lat, long, heading float64)
	Pipe(images chan GiftImage)
}

type GiftServer struct {
	db     *geoip2.Reader
	width  int
	height int
	source GiftImageSource
}

type GiftImage struct {
	img           *image.Paletted
	frameTimeMS   int
	disposalFlags uint8
	offset        image.Point
}

func NewGiftServer(w, h int, source GiftImageSource) GiftServer {
	gs := GiftServer{width: w, height: h, source: source}

	var err error
	gs.db, err = geoip2.Open("GeoLite2-City.mmdb")
	if err != nil {
		log.Fatal("Error opening geoip database: %+v", err)
	}

	return gs
}
