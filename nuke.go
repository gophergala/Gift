package gift

import (
	"image"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"log"
	"net/http"
	"os"
	"time"
)

type GiftImageNuke struct {
	lat, long     float64
	width, height int
	httpImages    chan GiftImage
}

func measure(fun func(), desc string) {
	start := time.Now()
	fun()
	end := time.Now()
	log.Printf("%s took: %+v", desc, end.Sub(start))
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

func overlayGif(path string, bounds image.Rectangle, images chan GiftImage) error {
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

	for i, frame := range frames.Image {

		parentWidth := bounds.Dx()
		parentHeight := bounds.Dy()
		childWidth := frame.Bounds().Dx()
		childHeight := frame.Bounds().Dy()

		offsetPt := image.Pt(parentWidth/2-childWidth/2, parentHeight/2-childHeight/2)

		images <- GiftImage{img: frame, frameTimeMS: frames.Delay[i], disposalFlags: DISPOSAL_RESTORE_PREV, offset: offsetPt}
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
		for i := 1; i < 7; i++ {
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

		measure(func() {
			overlayGif("nuke/explosion.gif", img.Bounds(), g.httpImages)
		}, "nuke overlay")

		black := image.NewPaletted(img.Bounds(), palette.Plan9)
		draw.Src.Draw(black, img.Bounds(), image.Black, image.Pt(0, 0))

		g.httpImages <- GiftImage{img: black, frameTimeMS: 2000}
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
