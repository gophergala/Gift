package gift

import (
	"code.google.com/p/freetype-go/freetype"
	"code.google.com/p/freetype-go/freetype/truetype"
	"github.com/oschwald/geoip2-golang"

	"fmt"
	"image"
	"image/color/palette"
	"io/ioutil"
	"log"
	"time"
)

type GiftImageCounter struct {
	font          *truetype.Font
	c             *freetype.Context
	city          *geoip2.City
	width, height int
}

func (g *GiftImageCounter) Geo(record *geoip2.City) {
	g.city = record
}

func (g *GiftImageCounter) Pipe(images chan *image.Paletted) {
	log.Printf("About to send GiftImageCounter")
	for i := 0; i < 10; i++ {
		img := image.NewPaletted(image.Rect(0, 0, g.width, g.height), palette.Plan9)
		g.c.SetDst(img)
		pt := freetype.Pt(g.width/2-100, g.height/2)
		g.c.DrawString(fmt.Sprintf("Frame: %d", i), pt)

		images <- img
		time.Sleep(100 * time.Millisecond)
	}
	close(images)
}
func (g *GiftImageCounter) Setup(width, height int) {
	fontBytes, err := ioutil.ReadFile("TimesNewRoman.ttf")
	if err != nil {
		log.Println(err)
		return
	}
	g.font, err = freetype.ParseFont(fontBytes)
	if err != nil {
		log.Println(err)
		return
	}

	g.width = width
	g.height = height

	fg := image.White
	g.c = freetype.NewContext()
	g.c.SetDPI(72)
	g.c.SetFont(g.font)
	g.c.SetFontSize(48)
	g.c.SetClip(image.Rect(0, 0, width, height))
	g.c.SetSrc(fg)
}
