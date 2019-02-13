package kindleland

import (
	"image"
	"image/draw"
	"log"
	"strings"

	"github.com/golang/freetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
)

type Page struct {
	ParagraphStart int
	ParagraphEnd   int
	RuneStart      int
	RuneEnd        int
}

type TextView struct {
	Bounds     image.Rectangle
	Text       string
	Paragraphs []string
	Pages      []Page
	Page       int
}

func NewTextView(text string, bounds image.Rectangle) *TextView {
	var paragraphs []string
	for _, p := range strings.Split(text, "\n") {
		if len(p) > 0 {
			paragraphs = append(paragraphs, p)
		}
	}

	return &TextView{
		Text:       text,
		Bounds:     bounds,
		Paragraphs: paragraphs,
	}
}

func (tv *TextView) Render() *image.RGBA {
	fg, bg := image.White, image.Black

	dpi := 168.0
	size := 12.0
	spacing := 1.5
	rgba := image.NewRGBA(tv.Bounds)

	f, err := freetype.ParseFont(goregular.TTF)
	if err != nil {
		log.Println(err)
		return rgba
	}

	draw.Draw(rgba, rgba.Bounds(), bg, image.ZP, draw.Src)

	c := freetype.NewContext()
	c.SetDPI(dpi)
	c.SetFont(f)
	c.SetFontSize(size)
	c.SetClip(tv.Bounds)
	c.SetDst(rgba)
	c.SetSrc(fg)
	c.SetHinting(font.HintingNone)

	scale := c.PointToFixed(12)

	min := freetype.Pt(tv.Bounds.Min.X, tv.Bounds.Min.Y+int(c.PointToFixed(size)>>6))
	pt := min
	max := freetype.Pt(tv.Bounds.Max.X, tv.Bounds.Max.Y)
	for _, p := range tv.Paragraphs {
		for _, r := range p {
			index := f.Index(r)
			hmetric := f.HMetric(scale, index)
			if pt.X+hmetric.AdvanceWidth > max.X {
				pt.Y += c.PointToFixed(size * spacing)
				pt.X = min.X
			}
			pt, err = c.DrawString(string(r), pt)
			if err != nil {
				log.Println(err)
				return rgba
			}
		}
		pt.Y += c.PointToFixed(size * spacing)
	}

	return rgba
}
