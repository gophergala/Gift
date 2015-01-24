package gift

import (
	"image"
	"image/gif"
	"log"
	"net/http"
	"os"
)

type GiftImageNuke struct {
	lat, long     float64
	width, height int
	httpImages    chan GiftImage
}

func (g *GiftImageNuke) Geo(lat, long, heading float64) {
	g.lat = lat
	g.long = long

	go func() {
		defer close(g.httpImages)
		launchFile, err := os.Open("nuke/nasr.gif")
		if err != nil {
			log.Printf("Unable to open intro nuke gif")
			return
		}
		frames, err := gif.DecodeAll(launchFile)
		if err != nil {
			log.Printf("Unable to decode intro nuke gif")
			return
		}
		for i := range frames.Image {
			g.httpImages <- GiftImage{img: frames.Image[i], frameTimeMS: frames.Delay[i]}
		}

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
			g.httpImages <- GiftImage{img: img.(*image.Paletted), frameTimeMS: 100}
		}
	}()
}

func (g *GiftImageNuke) Pipe(images chan GiftImage) {
	log.Printf("About to send nuke map")
	for pm := range g.httpImages {
		images <- pm
	}
	close(images)
}
func (g *GiftImageNuke) Setup(width, height int) {
	g.width = width
	g.height = height
	g.httpImages = make(chan GiftImage)
}
