package main

import (
	"fmt"
	"image"

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

	var text string

	for {
		kevent := <-channel
		fmt.Println(kevent)

		if kevent.Type == kindleland.KeyDown {
			fmt.Print(kevent.Name())

			text += kevent.Value()

			go func() {
				tv := kindleland.NewTextView(text, image.Rect(50, 50, 550, 750))
				tv.Size = 24

				img, err := tv.Render()
				if err != nil {
					panic(err)
				}

				fb.ApplyImage(img)

				if err := fb.UpdateScreenFx(kindleland.FxUpdateFast); err != nil {
					panic(err)
				}
			}()

			// say := exec.Command("say", letter)
			// if err := say.Run(); err != nil {
			// fmt.Println(err)
			// }
		}
	}
}
