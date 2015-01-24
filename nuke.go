package gift

import (
	"image"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"log"
	"math"
	"math/rand"
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

func overlayGif(path string, bounds image.Rectangle, images chan GiftImage) error {
	return embedGif(path, bounds, image.Pt(0, 0), DISPOSAL_RESTORE_BG, images)
}

func embedGif(path string, bounds image.Rectangle, offset image.Point, disposalFlags uint8, images chan GiftImage) error {
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

func embedFrame(path string, bounds image.Rectangle, offset image.Point, disposalFlags uint8, frameTimeMS int, images chan GiftImage) error {
	gifFile, err := os.Open(path)
	if err != nil {
		log.Printf("Unable to open gif: %s", path)
		return err
	}
	defer gifFile.Close()
	frame, err := gif.Decode(gifFile)
	if err != nil {
		log.Printf("Unable to decode gif: %s", path)
		return err
	}

	center := image.Pt(frame.Bounds().Dx()/2, frame.Bounds().Dy()/2)
	offsetPt := image.Pt(0, 0)
	offsetPt.X += offset.X - center.X
	offsetPt.Y += offset.Y - center.Y

	images <- GiftImage{img: frame.(*image.Paletted), frameTimeMS: frameTimeMS, disposalFlags: DISPOSAL_RESTORE_PREV, offset: offsetPt}
	return nil
}

func (g *GiftImageNuke) Geo(lat, long, heading float64) {
	g.lat = lat
	g.long = long

	go func() {
		defer close(g.httpImages)

		measure(func() {
			overlayGif("nuke/nasr.gif", image.Rect(0, 0, g.width, g.height), g.httpImages)
		}, "rocket launch image")

		var img image.Image
		measure(func() {
			for i := 1; i < 7; i++ {
				maptype := "roadmap"
				if i > 4 {
					maptype = "satellite"
				}

				url := mapURL(lat, long, g.width, g.height, i*3, maptype)
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
				g.httpImages <- GiftImage{img: img.(*image.Paletted), frameTimeMS: 10}
				center := image.Pt(img.Bounds().Dx()/2, img.Bounds().Dy()/2)
				startSide := rand.Intn(4)
				var startPt = image.Pt(0, 0)
				switch startSide {
				case 0: // Random Y coord, X = 0
					startPt.X = 0
					startPt.Y = rand.Intn(img.Bounds().Dy())
				case 1: // Random Y coord, X = right side
					startPt.X = img.Bounds().Dx()
					startPt.Y = rand.Intn(img.Bounds().Dy())
				case 2: // Random X coord, Y = 0
					startPt.X = rand.Intn(img.Bounds().Dx())
					startPt.Y = 0
				case 3:
					startPt.X = rand.Intn(img.Bounds().Dx())
					startPt.Y = img.Bounds().Dy()
				}
				delta := center.Sub(startPt)
				dlen := math.Sqrt(float64(delta.X*delta.X + delta.Y*delta.Y))
				dx, dy := float64(delta.X)/dlen, float64(delta.Y)/dlen

				crosshairSteps := 20
				timeStep := float64(10)
				for j := 0; j < crosshairSteps; j++ {

					t := (float64(j) / float64(crosshairSteps)) * dlen

					pt := image.Pt(startPt.X+int(dx*t), startPt.Y+int(dy*t))
					log.Printf("Point: %+v %f %f %f", pt, dx, dy, t)

					ts := int(timeStep)
					if j == crosshairSteps-1 {
						embedFrame("nuke/crosshair.gif", img.Bounds(), pt, DISPOSAL_RESTORE_BG, ts+100, g.httpImages)
					} else {
						embedFrame("nuke/crosshair_small.gif", img.Bounds(), pt, DISPOSAL_RESTORE_BG, ts, g.httpImages)

					}
				}
			}
		}, "google maps queries")

		measure(func() {
			overlayGif("nuke/explosion.gif", img.Bounds(), g.httpImages)
		}, "nuke overlay")

		black := image.NewPaletted(img.Bounds(), palette.Plan9)
		draw.Src.Draw(black, img.Bounds(), image.Black, image.Pt(0, 0))

		g.httpImages <- GiftImage{img: black, frameTimeMS: 200}
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
