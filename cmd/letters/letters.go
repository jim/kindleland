package main

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/fogleman/gg"

	"github.com/jim/kindleland"
)

func update(fb *kindleland.FrameBuffer, ctx *gg.Context, text string) error {
	ctx.Clear()
	ctx.SetRGB(0, 0, 0)
	ctx.Fill()

	ctx.SetRGB(1, 1, 1)
	ctx.DrawStringAnchored(text, 300, 390, .5, .5)
	fmt.Print("4 ")

	img := ctx.Image()
	fmt.Print("5 ")
	if err := fb.ApplyImage(img); err != nil {
		return err
	}
	fmt.Print("6 ")

	return fb.UpdateScreen()
}

func main() {
	// fb, err := kindleland.NewFrameBuffer("/dev/fb0", 600, 800)
	// if err != nil {
	// panic(err)
	// }

	// font, err := truetype.Parse(gobold.TTF)
	// if err != nil {
	// panic(err)
	// }
	// face := truetype.NewFace(font, &truetype.Options{Size: 450})

	channel, err := kindleland.NewKeyboardListener("/dev/input/event1")
	if err != nil {
		panic(err)
	}
	for {
		kevent := <-channel
		if kevent.Type == kindleland.KeyDown {
			fmt.Print(kevent.Name())

			letter := strings.ToLower(kevent.Name())
			go func() {
				// fmt.Println("1 ")
				// ctx := gg.NewContext(fb.Width, fb.Height)
				// fmt.Println("2 ")
				// ctx.SetFontFace(face)
				// fmt.Println("3 ")
				// update(fb, ctx, kevent.Name())
				// fmt.Println("7")
			}()
			say := exec.Command("say", letter)
			say.Run()
		}
	}
}
