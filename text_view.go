package kindleland

import (
	"fmt"
	"image"
	"image/draw"
	"log"

	"github.com/golang/freetype/truetype"
	"github.com/pkg/errors"

	"github.com/golang/freetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/math/fixed"
)

type Page struct {
	ParagraphStart int
	ParagraphEnd   int
	RuneStart      int
	RuneEnd        int
}

type TextView struct {
	Bounds image.Rectangle
	Text   string
	Pages  []Page
	Page   int
	Buffer *TextBuffer
	Size   float64
}

func NewTextView(text string, bounds image.Rectangle) *TextView {
	return &TextView{
		Text:   text,
		Bounds: bounds,
		Buffer: NewTextBuffer(text),
	}
}

func (tv *TextView) Render() (*image.RGBA, error) {
	fg, bg := image.White, image.Black

	dpi := 168.0
	spacing := 1.5
	rgba := image.NewRGBA(tv.Bounds)

	f, err := freetype.ParseFont(goregular.TTF)
	if err != nil {
		log.Println(err)
		return nil, errors.Wrap(err, "failed to parse font")
	}

	draw.Draw(rgba, rgba.Bounds(), bg, image.ZP, draw.Src)

	c := freetype.NewContext()
	c.SetDPI(dpi)
	c.SetFont(f)
	c.SetFontSize(tv.Size)
	c.SetClip(tv.Bounds)
	c.SetDst(rgba)
	c.SetSrc(fg)
	c.SetHinting(font.HintingFull)

	scale := c.PointToFixed(12)

	min := freetype.Pt(tv.Bounds.Min.X, tv.Bounds.Min.Y+int(c.PointToFixed(tv.Size)>>6))
	pt := min
	max := freetype.Pt(tv.Bounds.Max.X, tv.Bounds.Max.Y)
	words := 0

	for {
		word, space, ok := tv.Buffer.NextWord()
		fmt.Printf("%s, %s, %v\n", word, space, ok)
		if !ok {
			fmt.Println("not ok")
			break
		}
		width := wordWidth(word, scale, f)
		fmt.Println(word, width)
		if pt.X+width >= max.X {
			pt.Y += c.PointToFixed(tv.Size * spacing)
			pt.X = min.X
		}
		if pt.Y > max.Y {
			fmt.Println("hit bottom of view")
			break
		}
		words++
		pt, err = c.DrawString(word, pt)
		if err != nil {
			log.Println(err)
			return nil, errors.Wrap(err, "could not draw word")
		}
		pt, err = c.DrawString(space, pt)
		if err != nil {
			log.Println(err)
			return nil, errors.Wrap(err, "could not draw space")
		}
	}

	fmt.Printf("wrote %d words\n", words)

	return rgba, nil
}

func wordWidth(word string, scale fixed.Int26_6, font *truetype.Font) fixed.Int26_6 {
	sum := fixed.I(0)
	for _, r := range word {
		index := font.Index(r)
		hmetric := font.HMetric(scale, index)
		sum += hmetric.AdvanceWidth
	}
	return sum
}
