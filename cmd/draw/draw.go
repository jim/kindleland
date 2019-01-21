package main

import (
	"fmt"
	"math"
	"os"
	"syscall"
)

func NewFrameBuffer(device string) (*FrameBuffer, error) {
	file, err := os.OpenFile("/dev/fb0", os.O_RDWR, 0)
	size := 240000
	defer file.Close()
	if err != nil {
		panic(err)
	}
	fd := int(file.Fd())
	fb, err := syscall.Mmap(fd, 0, size, syscall.PROT_WRITE|syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		panic(err)
	}
	return &FrameBuffer{
		buffer: fb,
	}, nil
}

type FrameBuffer struct {
	buffer []byte
}

func (fb *FrameBuffer) Pixel(x, y, level int) error {
	if level < 0 || level > 15 {
		return fmt.Errorf("level must be between 0 and 15, got %d", level)
	}
	offset := x/2 + (y * 300)
	bits := int(fb.buffer[offset])

	var newBits int
	if x%2 == 0 {
		newBits = (bits & 15) + (level * 16)
	} else {
		newBits = (bits & 240) + level
	}

	fb.buffer[offset] = byte(newBits)
	return nil
}

func main() {
	fb, err := NewFrameBuffer("/dev/fb0")
	if err != nil {
		panic(err)
	}
	for x := 0; x < 600; x++ {
		for y := 0; y < 800; y++ {
			level := 15 - int(math.Abs(float64((y/50)-(x/38))))
			if err := fb.Pixel(x, y, level); err != nil {
				fmt.Printf("Failed at x: %d, y: %d, level: %d\n", x, y, level)
				panic(err)
			}
		}
	}
	// fmt.Scanln()
}
