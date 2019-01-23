package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math/rand"
	"os"
	"path"
	"strings"
)

func convertValue(value uint8, min, max int) uint8 {
	return uint8(min + int(float64(value)/255*float64(max-min)))
}

func addNoise(value uint8) uint8 {
	return value + 4 - uint8(rand.Intn(8))
}

func convert(c color.Color) color.Color {
	gray8 := c.(color.Gray)
	return color.Gray{
		Y: addNoise(convertValue(gray8.Y, 45, 212)),
	}
}

var filename = flag.String("input", "", "path to file to convert")

func main() {
	flag.Parse()
	if *filename == "" {
		panic("must specify a file as -input")
	}

	file, err := os.Open(*filename)
	defer file.Close()
	if err != nil {
		panic(err)
	}

	img, err := png.Decode(file)
	grayImg := img.(*image.Gray)
	if err != nil {
		panic(err)
	}
	output := image.NewGray(grayImg.Bounds())

	for y := grayImg.Bounds().Min.Y; y < grayImg.Bounds().Max.Y; y++ {
		for x := grayImg.Bounds().Min.X; x < grayImg.Bounds().Max.X; x++ {
			output.Set(x, y, convert(grayImg.At(x, y)))
		}
	}

	basename := strings.TrimSuffix(*filename, path.Ext(*filename))
	outputFile, err := os.Create(fmt.Sprintf("%s-einked.png", basename))
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()

	if err := png.Encode(outputFile, output); err != nil {
		panic(err)
	}
}
