package gift

import (
	"fmt"
	"image"
	"image/gif"
	"log"
	"net/http"
)

type GiftImageMap struct {
	lat, long     float64
	width, height int
	httpImages    chan GiftImage
}

var mapRoot = "https://maps.googleapis.com/maps/api/staticmap"

func mapURL(centerX, centerY float64, width, height, zoom int, maptype string) string {
	return fmt.Sprintf("%s?center=%f,%f&zoom=%d&size=%dx%d&format=gif&maptype=%s", mapRoot, centerX, centerY, zoom, width, height, maptype)
}

func (g *GiftImageMap) Geo(lat, long, heading float64) {
	g.lat = lat
	g.long = long

	go func() {
		defer close(g.httpImages)

		for i := 1; i < 6; i++ {
			url := mapURL(lat, long, g.width, g.height, i*3, "roadmap")
			resp, err := http.Get(url)
			if err != nil {
				log.Printf("Error requesting map: %d: %+v\n", i, err)
				continue
			}
			img, err := gif.Decode(resp.Body)
			if err != nil {
				log.Printf("Error decoding map: %+v", err)
				continue
			}
			g.httpImages <- GiftImage{img: img.(*image.Paletted), frameTimeMS: 100}
		}
	}()
}

func (g *GiftImageMap) Pipe(images chan GiftImage) {
	defer close(images)
	log.Printf("About to send map")
	for pm := range g.httpImages {
		images <- pm
	}
}
func (g *GiftImageMap) Setup(width, height int) {
	g.width = width
	g.height = height
	g.httpImages = make(chan GiftImage)
}
