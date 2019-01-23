package main

import (
	"fmt"
	"math"

	"github.com/jim/kindleland"
)

func main() {
	fb, err := kindleland.NewFrameBuffer("/dev/fb0", 600, 800)
	if err != nil {
		panic(err)
	}
	for x := 0; x < 600; x++ {
		for y := 0; y < 800; y++ {
			level := 15 - int(math.Abs(float64((y/50)-(x/38))))
			if err := fb.Pixel(x, y, uint8(level)); err != nil {
				fmt.Printf("Failed at x: %d, y: %d, level: %d\n", x, y, level)
				panic(err)
			}
		}
	}
	fb.UpdateScreen()
}
