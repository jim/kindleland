package main

import (
	"fmt"

	"github.com/jim/kindleland"
)

func main() {
	channel, err := kindleland.NewKeyboardListener("/dev/input/event0")
	if err != nil {
		panic(err)
	}
	for {
		kevent := <-channel
		fmt.Printf("%+v\n", kevent)
	}
}
