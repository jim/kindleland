package main

import (
	"image/color"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"github.com/jim/kindleland"
	"golang.org/x/image/font/gofont/goregular"
)

func main() {
	fb, err := kindleland.NewFrameBuffer("/dev/fb0", 600, 800)
	if err != nil {
		panic(err)
	}

	ctx := gg.NewContext(fb.Width, fb.Height)
	ctx.DrawCircle(300, 400, 100)
	ctx.SetRGB(1, 1, 1)
	ctx.Fill()

	ctx.SetRGB(0, 0, 0)
	font, err := truetype.Parse(goregular.TTF)
	if err != nil {
		panic(err)
	}
	face := truetype.NewFace(font, &truetype.Options{Size: 48})

	message := "Hello!"
	ctx.SetFontFace(face)
	ctx.DrawStringAnchored(message, 300, 395, .5, .5)

	ctx.SetStrokeStyle(gg.NewSolidPattern(color.White))

	for i := 0; i < 7; i++ {
		ctx.DrawCircle(300, 400, float64(100+i*15))
		ctx.SetLineWidth(7 - float64(i))
		ctx.Stroke()
	}

	err = fb.ApplyImage(ctx.Image())
	if err != nil {
		panic(err)
	}

	err = fb.UpdateScreen()
	if err != nil {
		panic(err)
	}
}
