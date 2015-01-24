package gift

import (
	"image"
	"image/draw"
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

func embedGif(path string, images chan GiftImage) error {
	gifFile, err := os.Open(path)
	if err != nil {
		log.Printf("Unable to open gif: %s", path)
		return err
	}
	defer gifFile.Close()
	frames, err := gif.DecodeAll(gifFile)
	if err != nil {
		log.Printf("Unable to decode gif: %s", path)
		return err
	}
	for i := range frames.Image {
		images <- GiftImage{img: frames.Image[i], frameTimeMS: frames.Delay[i]}
	}
	return nil
}

func overlayGif(path string, img *image.Paletted, images chan GiftImage) error {
	gifFile, err := os.Open(path)
	if err != nil {
		log.Printf("Unable to open gif: %s", path)
		return err
	}
	defer gifFile.Close()
	frames, err := gif.DecodeAll(gifFile)
	if err != nil {
		log.Printf("Unable to decode gif: %s", path)
		return err
	}
	for i := range frames.Image {
		output := image.NewPaletted(img.Bounds(), img.Palette)
		draw.Over.Draw(output, img.Bounds(), frames.Image[i], image.Pt(0, 0))
		images <- GiftImage{img: output, frameTimeMS: frames.Delay[i]}
	}
	return nil
}

func (g *GiftImageNuke) Geo(lat, long, heading float64) {
	g.lat = lat
	g.long = long

	go func() {
		defer close(g.httpImages)

		embedGif("nuke/nasr.gif", g.httpImages)

		var img image.Image
		for i := 0; i < 8; i++ {
			url := mapURL(lat, long, g.width, g.height, i*3)
			resp, err := http.Get(url)
			if err != nil {
				log.Printf("Error requesting map: %d: %+v\n", i, err)
				continue
			}
			// Reuse our image we declared outside this loop so we can
			// overlay on top of this last frame
			img, err = gif.Decode(resp.Body)
			if err != nil {
				log.Printf("Error decoding map: %+v", err)
				continue
			}
			g.httpImages <- GiftImage{img: img.(*image.Paletted), frameTimeMS: 100}
		}

		overlayGif("nuke/explosion.gif", img.(*image.Paletted), g.httpImages)
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
