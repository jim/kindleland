package main

import (
	"os/exec"
	"strings"

	"github.com/jim/kindleland"
)

func main() {
	keyboard, err := kindleland.NewKeyboardListener("/dev/input/event0")
	if err != nil {
		panic(err)
	}
	fiveWay, err := kindleland.NewKeyboardListener("/dev/input/event1")
	if err != nil {
		panic(err)
	}
	channel := make(chan kindleland.KeyboardEvent)
	go func() {
		for {
			select {
			case ke := <-keyboard:
				channel <- ke
			case ke := <-fiveWay:
				channel <- ke
			}
		}
	}()
	for {
		kevent := <-channel
		if kevent.Type == kindleland.KeyDown {
			letter := strings.ToLower(kevent.Name())
			say := exec.Command("say", letter)
			say.Run()
		}
	}
}
