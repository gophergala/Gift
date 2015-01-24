package gift

import (
	"code.google.com/p/freetype-go/freetype"
	"code.google.com/p/freetype-go/freetype/truetype"
	"github.com/oschwald/geoip2-golang"

	"image"
	"log"
)

type GiftImageCounter struct {
	font          *truetype.Font
	c             *freetype.Context
	width, height int
}

type GiftImageSource interface {
	Setup(width, height int)
	Pipe(images chan *image.Paletted, record *geoip2.City)
}

type GiftServer struct {
	db     *geoip2.Reader
	width  int
	height int
	source GiftImageSource
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
