package kindleland

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"syscall"
)

func Gray4Downsample(c color.Color) uint8 {
	red, _, _, _ := c.RGBA()
	return uint8(red) >> 4
}

func NewFrameBuffer(device string, width, height int) (*FrameBuffer, error) {
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
		Width:  width,
		Height: height,
	}, nil
}

func (fb *FrameBuffer) Pixel(x, y int, level uint8) error {
	if level < 0 || level > 15 {
		return fmt.Errorf("level must be between 0 and 15, got %d", level)
	}
	offset := x/2 + (y * fb.Width / 2)
	if offset >= len(fb.buffer) {
		return fmt.Errorf("%d is out of range; max is %d; x: %d, y: %d", offset, len(fb.buffer)-1, x, y)
	}

	bits := uint8(fb.buffer[offset])

	var newBits uint8
	if x%2 == 0 {
		newBits = (bits & 15) + (level * 16)
	} else {
		newBits = (bits & 240) + level
	}

	fb.buffer[offset] = byte(newBits)
	return nil
}

type FrameBuffer struct {
	buffer []byte
	Width  int
	Height int
}

func (fb *FrameBuffer) ApplyImage(img image.Image) error {
	for y := 0; y < fb.Height; y++ {
		for x := 0; x < fb.Width; x++ {
			color := img.At(x, y)
			gray := Gray4Downsample(color)
			err := fb.Pixel(x, y, gray)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Flush any changes to the framebuffer to the display
func (fb *FrameBuffer) UpdateScreen() error {
	file, err := os.OpenFile("/proc/eink_fb/update_display", os.O_WRONLY, 0)
	defer file.Close()
	if err != nil {
		return err
	}
	if _, err := file.Write([]byte("1\n")); err != nil {
		return err
	}
	return nil
}
