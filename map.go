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
	httpImages    chan *image.Paletted
}

var mapRoot = "https://maps.googleapis.com/maps/api/staticmap"

func mapURL(centerX, centerY float64, width, height, zoom int) string {
	return fmt.Sprintf("%s?center=%f,%f&zoom=%d&size=%dx%d&format=gif", mapRoot, centerX, centerY, zoom, width, height)
}

func (g *GiftImageMap) Geo(lat, long, heading float64) {
	g.lat = lat
	g.long = long

	go func() {

		for i := 0; i < 8; i++ {
			url := mapURL(lat, long, g.width, g.height, i*3)
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
			g.httpImages <- img.(*image.Paletted)
		}
		close(g.httpImages)
	}()
}

func (g *GiftImageMap) Pipe(images chan *image.Paletted) {
	log.Printf("About to send map")
	for pm := range g.httpImages {
		images <- pm
	}
	close(images)
}
func (g *GiftImageMap) Setup(width, height int) {
	g.width = width
	g.height = height
	g.httpImages = make(chan *image.Paletted)
}
