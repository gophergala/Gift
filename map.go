package gift

import (
	"github.com/oschwald/geoip2-golang"

	"fmt"
	"image"
	"image/gif"
	"log"
	"net/http"
)

type GiftImageMap struct {
	city          *geoip2.City
	width, height int
	httpImages    chan *image.Paletted
}

var mapRoot = "https://maps.googleapis.com/maps/api/staticmap"

func mapURL(centerX, centerY float64, width, height, zoom int) string {
	return fmt.Sprintf("%s?center=%f,%f&zoom=%d&size=%dx%d&format=gif", mapRoot, centerX, centerY, zoom, width, height)
}

func (g *GiftImageMap) Geo(record *geoip2.City) {
	g.city = record

	go func() {

		for i := 0; i < 8; i++ {
			url := mapURL(record.Location.Latitude, record.Location.Longitude, g.width, g.height, i*3)
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