package main

import (
	"fmt"
	"image"
	"os/exec"
	"strings"

	"github.com/jim/kindleland"
)

func main() {
	fb, err := kindleland.NewFrameBuffer("/dev/fb0", 600, 800)
	if err != nil {
		panic(err)
	}

	channel, err := kindleland.NewKeyboardListener("/dev/input/event0")
	if err != nil {
		panic(err)
	}
	for {
		kevent := <-channel
		fmt.Println(kevent)

		if kevent.Type == kindleland.KeyDown {
			fmt.Print(kevent.Name())

			letter := strings.ToLower(kevent.Name())
			go func() {
				tv := kindleland.NewTextView(letter, image.Rect(50, 50, 550, 750))

				fb.ApplyImage(tv.Render())

				if err := fb.UpdateScreen(); err != nil {
					panic(err)
				}
			}()

			say := exec.Command("say", letter)
			if err := say.Run(); err != nil {
				fmt.Println(err)
			}
		}
	}
}
