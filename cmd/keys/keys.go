package main

import (
	"fmt"
	"os"
)

func main() {
	keyboard, err := os.Open("/dev/input/event0")
	defer keyboard.Close()
	if err != nil {
		panic(err)
	}
	buf := make([]byte, 16)

	for {
		n, err := keyboard.Read(buf)
		if err != nil {
			panic(err)
		}
		fmt.Println(n, buf)
	}
}
